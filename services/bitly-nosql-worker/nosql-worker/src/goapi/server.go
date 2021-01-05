package main

import (
	"database/sql"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/codegangsta/negroni"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
	"github.com/unrolled/render"
)

//MySQL Connection
// var mysql_connect = "cmpe281:cmpe281@tcp(10.0.5.117:3306)/cmpe281"
var mysql_connect = "root:cmpe281@tcp(mysql:3306)/cmpe281"

//NoSQL Connection
// nosql_connect = "http://internal-bilty-nosql-clb-1559610795.us-west-2.elb.amazonaws.com:9090"
var nosql_connect = "http://api_node_1:9090"

// RabbitMQ Config
// var rabbitmq_server = "10.0.5.243"
var rabbitmq_server = "rabbitmq"
var rabbitmq_port = "5672"
var rabbitmq_user = "guest"
var rabbitmq_pass = "guest"

// Queue Array
var queueList []string

// NewServer configures and returns a Server.
func NewServer() *negroni.Negroni {
	formatter := render.New(render.Options{
		IndentJSON: true,
	})
	n := negroni.Classic()
	mx := mux.NewRouter()
	initRoutes(mx, formatter)
	n.UseHandler(mx)
	queueList = append(queueList, "shortlinks_used", "new_shortlink")
	initQueueConsumer(queueList)
	return n
}

// API Routes
func initRoutes(mx *mux.Router, formatter *render.Render) {
	mx.HandleFunc("/ping", pingHandler(formatter)).Methods("GET")
}

// Helper Functions
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

// API Ping Handler
func pingHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		formatter.JSON(w, http.StatusOK, struct{ Test string }{"NoSQL Worker alive!"})
	}
}

// Queue Consumer Initialization
func initQueueConsumer(queueList []string) {
	for _, q := range queueList {
		c := make(chan LinkTrendData)
		consumerRoutineStarter(q, c)
		switch q {
		case "new_shortlink":
			fmt.Println("Starting Link create ROUTINE")
			go createLinkRoutine(q, c)
		case "shortlinks_used":
			fmt.Println("Starting Link Hit ROUTINE")
			go updateHitRoutine(q, c)
		}
	}
}

// Start Consumer routine
func consumerRoutineStarter(queueName string, c chan LinkTrendData) {
	fmt.Println("Messanger CONSUME()" + queueName)
	go queue_receive(queueName, c)
}

// updateHit routine
func updateHitRoutine(queue string, c chan LinkTrendData) {
	fmt.Println("Inside Link Hit routine")
	for linkTrendData := range c {
		fmt.Println(linkTrendData)
		slice := strings.Split(linkTrendData.ShortLink, "/")
		key := slice[len(slice)-1]
		fmt.Println("key: " + key)
		resp, err := http.Get(nosql_connect + "/api/" + key)
		if err != nil {
			fmt.Println("createLinkRoutine(): nosql get failed" + err.Error())
			//c <- linkTrendData
			if resp != nil {
				fmt.Println("createLinkRoutine(): nosql get failed", resp.Body)
				resp.Body.Close()
			}
			continue
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			//c <- linkTrendData
			if resp != nil {
				resp.Body.Close()
			}
			continue
		}
		err = json.Unmarshal(body, &linkTrendData)
		linkTrendData.Count++
		jsonData, err := json.Marshal(linkTrendData)
		client := &http.Client{}
		fmt.Println("Update the count in the DB: ", linkTrendData.Count)

		// Update in MySQL DB
		db, err := sql.Open("mysql", mysql_connect)
		defer db.Close()
		if err != nil {
			fmt.Println(err)
			continue
		}
		update, err := db.Prepare("UPDATE linkdetails set hitcount=? WHERE shortlink like '%" + key + "'")
		if err != nil {
			fmt.Println(err)
			continue
		}
		update.Exec(linkTrendData.Count)
		defer update.Close()

		// Update in NoSQL DB
		request, err := http.NewRequest(http.MethodPut, nosql_connect + "/api/"+key, bytes.NewBuffer(jsonData)) //TODO : add actual load balancer address
		if err != nil {
			fmt.Println("createLinkRoutine(): nosql put failed" + err.Error())
			//c <- linkTrendData
			continue
		}
		request.Header.Set("Content-Type", "application/json; charset=utf-8")
		resp, err = client.Do(request)
		if (err != nil) || (resp.StatusCode != http.StatusOK) {
			//c <- linkTrendData
			if resp != nil {
				resp.Body.Close()
			}
			continue
		}

	}
}

// createLink routine
func createLinkRoutine(queue string, c chan LinkTrendData) {
	for linkTrendData := range c {
		fmt.Println(linkTrendData)
		jsonData, err := json.Marshal(linkTrendData)
		slice := strings.Split(linkTrendData.ShortLink, "/")
		key := slice[len(slice)-1]
		fmt.Println("key: " + key)

		// Insert into MySQL database
		db, err := sql.Open("mysql", mysql_connect)
		defer db.Close()
		if err != nil {
			log.Fatal(err)
		}
		insert, err := db.Prepare("INSERT INTO linkdetails (shortlink, uri, hitcount) VALUES ( ?,?,? )")
		if err != nil {
			log.Fatal(err)
		}
		insert.Exec(linkTrendData.ShortLink, linkTrendData.Uri, linkTrendData.Count)
		defer insert.Close()

		// Insert into NoSQL
		resp, err := http.Post(nosql_connect + "/api/"+key, "application/json", bytes.NewBuffer(jsonData)) //TODO : add actual load balancer address
		if (err != nil) || (resp.StatusCode != http.StatusOK) {
			fmt.Println("createLinkRoutine(): nosql post failed" + err.Error())
			continue
		}
	}
}

// Receive message from Queues to Process
func queue_receive(rabbitmq_queue string, c chan LinkTrendData) {
	conn, err := amqp.Dial("amqp://" + rabbitmq_user + ":" + rabbitmq_pass + "@" + rabbitmq_server + ":" + rabbitmq_port + "/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	fmt.Println("Messanger consumerRoutine(): " + rabbitmq_queue)
	linkTrendData := LinkTrendData{}

	q, err := ch.QueueDeclare(
		rabbitmq_queue, // name
		false,          // durable
		false,          // delete when usused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name,   // queue
		"orders", // consumer
		true,     // auto-ack
		false,    // exclusive
		false,    // no-local
		false,    // no-wait
		nil,      // args
	)
	failOnError(err, "Failed to register a consumer")

	// process queue messages
	for msg := range msgs {
		json.Unmarshal(msg.Body, &linkTrendData)
		c <- linkTrendData
	}
}

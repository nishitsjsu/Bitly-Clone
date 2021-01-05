package main

import (
	"database/sql"
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

// NoSQL Connection
var nosql_connect = "http://api_node_1:9090"

// RabbitMQ Config
// var rabbitmq_server = "10.0.5.243"
var rabbitmq_server = "rabbitmq"
var rabbitmq_port = "5672"
var rabbitmq_queue = "shortlinks_used"
var rabbitmq_user = "guest"
var rabbitmq_pass = "guest"

// NewServer configures and returns a Server.
func NewServer() *negroni.Negroni {
	formatter := render.New(render.Options{
		IndentJSON: true,
	})
	n := negroni.Classic()
	mx := mux.NewRouter()
	initRoutes(mx, formatter)
	n.UseHandler(mx)
	return n
}

// API Routes
func initRoutes(mx *mux.Router, formatter *render.Render) {
	mx.HandleFunc("/ping", pingHandler(formatter)).Methods("GET", "OPTIONS")
	mx.HandleFunc("/{shortLink}", linkRedirect(formatter)).Methods("GET", "OPTIONS")
	mx.HandleFunc("/all/trends", allTrendsHandler(formatter)).Methods("GET", "OPTIONS")
	mx.HandleFunc("/trends/{shortLink}", trendsHandler(formatter)).Methods("GET", "OPTIONS")
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
		formatter.JSON(w, http.StatusOK, struct{ Test string }{"Link redirect server alive!"})
	}
}

// Redirect 
func linkRedirect(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
		var linkData LinkData
		shortLink := "http://" + req.Host + req.URL.Path
		fmt.Println("RedirectLink(): ", linkData, shortLink)

		// NoSQL Get API
		fmt.Printf("fetching short link info from the nosql database")
		slice := strings.Split(shortLink, "/")
		key := slice[len(slice)-1]
		fmt.Println("key: " + key)
		resp, err := http.Get(nosql_connect + "/api/" + key)
		if err != nil {
			log.Fatalf("linkRedirect(): nosql get failed: " + err.Error())
			if resp != nil {
				fmt.Println("linkRedirect(): nosql get failed", resp.Body)
				resp.Body.Close()
			}
		}
		// Read data
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("linkRedirect(): error in reading the data: " + err.Error())
			if resp != nil {
				resp.Body.Close()
			}
		}
		err = json.Unmarshal(body, &linkData)

		
		fmt.Println("Redirecting user to this link: %v & Short link is %v ", linkData.Uri, linkData.ShortLink)
		// Publish to the queue
		queue_send(linkData)
		// Redirect
		http.Redirect(w, req, linkData.Uri, http.StatusSeeOther)
	}
}

// Trends handler
func trendsHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
		// var linkTrend LinkTrendData
		params := mux.Vars(req)
		shortLink := params["shortLink"]
		fmt.Println("Short Link is : ", shortLink)

		r, err := http.Get(nosql_connect + "/api/" + shortLink)
		if err != nil {
			if r != nil {
				defer r.Body.Close()
			}
			log.Fatalf("Error in the GET API " + err.Error())
		}
		fmt.Println("before readiing response body")
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			defer r.Body.Close()
			log.Fatalf("Error in parsing the data from GET API " + err.Error())
		}
		var responseTrendData LinkData
		err = json.Unmarshal(body, &responseTrendData)
		formatter.JSON(w, http.StatusOK, responseTrendData)
	}
}

// Get all links
func allTrendsHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
		// Get from MySQL
		db, err := sql.Open("mysql", mysql_connect)
		defer db.Close()
		if err != nil {
			log.Fatal(err)
		}

		query, err := db.Query("SELECT * FROM linkdetails")
		if err != nil {
			log.Fatal(err)
			defer query.Close()
		}
		defer query.Close()

		var linkDataArray []LinkData
		linkData := LinkData{}
		// var linkData LinkData
		for query.Next() {
			err = query.Scan(&linkData.Id, &linkData.ShortLink, &linkData.Uri, &linkData.Count)
			if err != nil {
				defer query.Close()
			}
			linkDataArray = append(linkDataArray, linkData)
		}

		// responseData, err := json.Marshal(linkDataArray)
		formatter.JSON(w, http.StatusOK, linkDataArray)
	}
}

// Send Order to Queue for Processing
func queue_send(linkData LinkData) {
	conn, err := amqp.Dial("amqp://" + rabbitmq_user + ":" + rabbitmq_pass + "@" + rabbitmq_server + ":" + rabbitmq_port + "/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		rabbitmq_queue, // name
		false,          // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// body := message
	jsonBody, err := json.Marshal(linkData)
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonBody,
		})
	log.Printf(" [x] is Sent")
	failOnError(err, "Failed to publish a message")
}

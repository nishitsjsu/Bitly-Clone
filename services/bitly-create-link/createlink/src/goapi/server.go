/*
	Create Link API in Go (Version 1)
	Uses MySQL
	Uses RabbitMQ
*/

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/codegangsta/negroni"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
	"github.com/streadway/amqp"
	"github.com/unrolled/render"
)

//Kong URI
var kong_uri = "Kong API/LRS/"

// RabbitMQ Config
var rabbitmq_server = "rabbitmq"
// var rabbitmq_server = "10.0.5.243"
var rabbitmq_port = "5672"
var rabbitmq_queue = "new_shortlink"
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
	mx.HandleFunc("/createlink", createLink(formatter)).Methods("POST", "OPTIONS")
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
		formatter.JSON(w, http.StatusOK, struct{ Test string }{"Create link server alive!"})
	}
}

//Create new link
func createLink(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
		var linkData LinkData
		err := json.NewDecoder(req.Body).Decode(&linkData)
		if err != nil {
			log.Fatal(err)
		}
		// Create Shortlink
		uuid := uuid.NewV4()
		// linkData.ShortLink = "Kong API/LRS/" + uuid.String()[:5]
		linkData.ShortLink = kong_uri + uuid.String()[:5]
		fmt.Println("linkData.ShortLink: %s, URI: %s ", linkData.ShortLink, linkData.Uri)
		fmt.Println("Link created:", linkData)
		// Publish to the queue
		queue_send(linkData)
		formatter.JSON(w, http.StatusOK, linkData)
	}
}

// Send linkdata to Queue for Processing
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
	jsonMessage, err := json.Marshal(linkData)
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonMessage,
		})
	log.Printf(" [x] Sent ")
	failOnError(err, "Failed to publish a message")
}
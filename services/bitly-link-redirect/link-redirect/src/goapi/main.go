/*
	Bitly link-redirect service in Go (Version 1)
	Uses MySQL, NoSQL database APIs & RabbitMQ
	Runs on port 3002
*/
package main

import (
	"os"
)

func main() {

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3002"
	}

	server := NewServer()
	server.Run(":" + port)
}

/*
	Bitly nosql worker service in Go (Version 1)
	Uses MySQL, NoSQL database APIs & RabbitMQ
	Runs on port 3004
*/

package main

import (
	"os"
)

func main() {

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3004"
	}

	server := NewServer()
	server.Run(":" + port)
}

/*
	Bitly Create-Link service in Go (Version 1)
	Uses RabbitMQ
	Runs on port 3001
*/

package main

import (
	"os"
)

func main() {

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3001"
	}

	server := NewServer()
	server.Run(":" + port)
}

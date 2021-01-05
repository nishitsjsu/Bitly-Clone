/*
	Bitly Create-Link service in Go (Version 1)
	Uses RabbitMQ
	Runs on port 3001
*/

package main

type LinkData struct {
	ShortLink string
	Uri       string
	Count     int
}

/*
	Bitly nosql worker service in Go (Version 1)
	Uses MySQL, NoSQL database APIs & RabbitMQ
	Runs on port 3004
*/

package main

type LinkTrendData struct {
	ShortLink string `json: "shortlink"`
	Uri       string `json: "uri"`
	Count     int    `json: "hits"`
}

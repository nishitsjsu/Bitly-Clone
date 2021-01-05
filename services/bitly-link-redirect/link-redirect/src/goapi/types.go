/*
	Bitly link-redirect service in Go (Version 1)
	Uses MySQL, NoSQL database APIs & RabbitMQ
	Runs on port 3002
*/
package main

type LinkData struct {
	Id        int    
	ShortLink string 
	Uri       string 
	Count int
}
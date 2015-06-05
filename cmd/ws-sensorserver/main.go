package main

import (
	"github.com/gonium/websocketsensors"
	"log"
	"net/http"
)

func main() {
	log.Printf("WS-Sensor Server started")
	http.HandleFunc("/echo", websocketsensors.EchoHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("Error: " + err.Error())
	}
}

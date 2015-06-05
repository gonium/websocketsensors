package main

import (
	"flag"
	"fmt"
	"github.com/gonium/websocketsensors"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const (
	defaultport = 8080
	defaulthost = "localhost"
)

var port = flag.Int("port", defaultport, "which port to listen on")
var host = flag.String("host", defaulthost, "which host to listen on")

func init() {
	flag.Parse()
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/echo", websocketsensors.EchoHandler)
	wssh, err := websocketsensors.NewWSSHandler()
	if err != nil {
		log.Fatal("Failed to create websocket sensor handler: %s",
			err.Error())
	}
	r.HandleFunc("/api/v0/sensor/update/{uuid}/{auth}", wssh.SensorUpdateHandler)
	r.HandleFunc("/api/v0/sensor/subscribe/{uuid}/{auth}",
		wssh.SensorSubscribeHandler)
	http.Handle("/", r)
	listenAddress := fmt.Sprintf("%s:%d", *host, *port)
	log.Println("Starting server at " + listenAddress)
	err = http.ListenAndServe(listenAddress, nil)
	if err != nil {
		panic("Failed to start http server: " + err.Error())
	}
}

package main

import (
	"flag"
	"fmt"
	"github.com/gonium/websocketsensors"
	"log"
	"time"
)

const (
	defaultaction = "echo"
	defaulthost   = "ws://localhost:8080"
)

var action = flag.String("action", defaultaction, "which subcommand to run, {echo|...}")
var host = flag.String("url", defaulthost, "endpoint of the websocket server")

func init() {
	flag.Parse()
	if *action == defaultaction {
		log.Printf("running default echo action (see -action parameter)\r\n")
	}
	if *host == defaulthost {
		log.Printf("Using default server at %s\r\n", *host)
	}
}

func main() {
	log.Printf("WS-Sensor started")
	var hosturl string
	if *action == "echo" {
		hosturl = fmt.Sprintf("%s/echo", *host)
	}
	client, err := websocketsensors.NewWSSClient(hosturl)
	if err != nil {
		log.Fatalf("Failed to connect to %s: %v", hosturl, err)
	}
	defer client.Close()

	switch *action {
	case "echo":
		for i := 0; i < 10; i++ {
			request := fmt.Sprintf("Hello %d", i)
			msg, err := client.SendRecv(request)
			if err != nil {
				log.Printf("Failed to send message: %s", err)
			} else {
				log.Printf("Requested %s, received answer: %s", request, msg)
			}
			time.Sleep(time.Second)
		}
		log.Printf("WS-Sensor finished")
	default:
		log.Fatalf("Unknown action %s - aborting\r\n", *action)
	}
}

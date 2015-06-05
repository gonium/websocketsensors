package main

import (
	"encoding/json"
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
var uuid = flag.String("uuid", "", "uuid of the sensor")
var auth = flag.String("auth", "", "auth token of the sensor")

func abortIfEmpty(param string, msg string) {
	if param == "" {
		log.Fatal(msg)
	}
}

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
	switch *action {

	case "echo":
		hosturl = fmt.Sprintf("%s/echo", *host)
		client, err := websocketsensors.NewWSSClient(hosturl)
		if err != nil {
			log.Fatalf("Failed to connect to %s: %v", hosturl, err)
		}
		defer client.Close()

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

	case "update":
		abortIfEmpty(*uuid, "Please provide a valid uuid")
		abortIfEmpty(*auth, "Please provide a valid authentication token")
		hosturl = fmt.Sprintf("%s/api/v0/sensor/update/%s/%s", *host, *uuid,
			*auth)
		client, err := websocketsensors.NewWSSClient(hosturl)
		if err != nil {
			log.Fatalf("Failed to connect to %s: %v", hosturl, err)
		}
		defer client.Close()
		for i := 0; i < 10; i++ {
			update := websocketsensors.UpdateMessage{Timestamp: 4711, Value: 23.42}
			request, err := json.Marshal(update)
			if err != nil {
				log.Fatal("Failed to marshal update: %s", err.Error())
			}
			rawresp, err := client.SendRecv(string(request))
			if err != nil {
				log.Printf("Failed to send message: %s", err)
			} else {
				log.Printf("send update: %s, raw response was %s", request,
					rawresp)
				var response websocketsensors.UpdateResponse
				err = json.Unmarshal([]byte(rawresp), &response)
				if err != nil {
					log.Printf("Failed to decode server response: %s\r\n",
						err.Error())
				} else {
					if response.Status == websocketsensors.STATUS_OK {
						log.Printf("Successfully stored update.")
					}
				}
			}
			time.Sleep(time.Second)
		}

	case "subscribe":
		abortIfEmpty(*uuid, "Please provide a valid uuid")
		abortIfEmpty(*auth, "Please provide a valid authentication token")
		hosturl = fmt.Sprintf("%s/api/v0/sensor/subscribe/%s/%s", *host, *uuid,
			*auth)
		client, err := websocketsensors.NewWSSClient(hosturl)
		if err != nil {
			log.Fatalf("Failed to connect to %s: %v", hosturl, err.Error())
		}
		defer client.Close()
		for {
			rawresp, err := client.Recv()
			if err != nil {
				log.Printf("Failed to receive message: %s", err.Error())
			} else {
				log.Printf("raw response was %s", rawresp)
				var response websocketsensors.UpdateMessage
				err = json.Unmarshal([]byte(rawresp), &response)
				if err != nil {
					log.Printf("Failed to decode server response: %s\r\n",
						err.Error())
				} else {
					log.Printf("Received: %d - %f\r\n", response.Timestamp,
						response.Value)
				}
			}
			time.Sleep(time.Second)
		}

	default:
		log.Fatalf("Unknown action %s - aborting\r\n", *action)
	}
	log.Printf("WS-Sensor finished")
}

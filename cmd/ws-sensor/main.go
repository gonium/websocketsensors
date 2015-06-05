package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

var cstDialer = websocket.Dialer{
	Subprotocols:    []string{"p1", "p2"},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func sendRecv(ws *websocket.Conn, message string) (response string, err error) {
	if err = ws.SetWriteDeadline(time.Now().Add(time.Second)); err != nil {
		return
	}
	if err = ws.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		return
	}
	if err = ws.SetReadDeadline(time.Now().Add(time.Second)); err != nil {
		return
	}
	_, byteresponse, readerr := ws.ReadMessage()
	if readerr != nil {
		return "", readerr
	}
	//log.Printf("message=%s, want %s", p, message)
	return string(byteresponse), nil
}

func main() {
	log.Printf("WS-Sensor started")
	ws, _, err := cstDialer.Dial("ws://localhost:8080/echo", nil)
	if err != nil {
		log.Fatalf("Dial: %v", err)
	}
	defer ws.Close()
	for i := 0; i < 10; i++ {
		msg, err := sendRecv(ws, fmt.Sprintf("Hello %d", i))
		if err != nil {
			log.Printf("Failed to send message: %s", err)
		} else {
			log.Printf("Received answer: %s", msg)
		}
		time.Sleep(time.Second)
	}
	log.Printf("WS-Sensor finished")

}

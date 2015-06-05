package websocketsensors

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

// Simple handler that just echoes whatever you send to it. Only for
// testing purposes.
func EchoHandler(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			return
		}
		err = conn.WriteMessage(messageType, p)
		if err != nil {
			return
		}
		log.Printf("Echoed: %s", string(p))
	}
}

type WSSHandler struct {
	upgrader websocket.Upgrader
}

func NewWSSHandler() (handler *WSSHandler, err error) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	handler = &WSSHandler{upgrader: upgrader}
	return handler, nil
}

func (wssh *WSSHandler) SensorUpdateHandler(w http.ResponseWriter, r *http.Request) {
	uuid := mux.Vars(r)["uuid"]
	auth := mux.Vars(r)["auth"]
	log.Printf("New update stream from %s (%s)\r\n", uuid, auth)
	conn, err := wssh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	conn.SetPingHandler(nil)
	conn.SetPongHandler(nil)
	for {
		//messageType, p, err := conn.ReadMessage()
		var msg UpdateMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			// probably the connection was closed
			return
		}
		log.Printf("Update received: %d - %f", msg.Timestamp, msg.Value)
		err = conn.WriteJSON(UpdateResponse{Status: STATUS_OK})
		if err != nil {
			return
		}
	}
}

func (wssh *WSSHandler) SensorSubscribeHandler(w http.ResponseWriter, r *http.Request) {
	uuid := mux.Vars(r)["uuid"]
	auth := mux.Vars(r)["auth"]
	log.Printf("New subscription stream for %s (%s)\r\n", uuid, auth)
	conn, err := wssh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	conn.SetPingHandler(nil)
	conn.SetPongHandler(nil)
	for {
		// TODO: Use hub to manage broadcast
		//messageType, p, err := conn.ReadMessage()
		var msg UpdateMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			// probably the connection was closed
			return
		}
		log.Printf("Update received: %d - %f", msg.Timestamp, msg.Value)
		err = conn.WriteJSON(UpdateResponse{Status: STATUS_OK})
		if err != nil {
			return
		}
	}
}

package websocketsensors

import (
	"github.com/gorilla/websocket"
	"time"
)

type WSSClient struct {
	wssDialer websocket.Dialer
	ws        *websocket.Conn
}

func NewWSSClient(host string) (client *WSSClient, err error) {
	dialer := websocket.Dialer{
		Subprotocols:    []string{"p1", "p2"},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024}
	websocket, _, err := dialer.Dial(host, nil)
	if err != nil {
		return nil, err
	}
	client = &WSSClient{wssDialer: dialer, ws: websocket}
	return client, nil
}

func (wss *WSSClient) Close() {
	wss.ws.Close()
}

func (wss *WSSClient) SendRecv(message string) (response string, err error) {
	if err = wss.ws.SetWriteDeadline(time.Now().Add(time.Second)); err != nil {
		return
	}
	if err = wss.ws.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		return
	}
	if err = wss.ws.SetReadDeadline(time.Now().Add(time.Second)); err != nil {
		return
	}
	_, byteresponse, readerr := wss.ws.ReadMessage()
	if readerr != nil {
		return "", readerr
	}
	//log.Printf("message=%s, want %s", p, message)
	return string(byteresponse), nil
}

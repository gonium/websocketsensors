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
	websocket.SetPingHandler(nil)
	websocket.SetPongHandler(nil)
	client = &WSSClient{wssDialer: dialer, ws: websocket}
	return client, nil
}

func (wss *WSSClient) Close() {
	wss.ws.Close()
}

func (wss *WSSClient) Send(message string) (err error) {
	if err = wss.ws.SetWriteDeadline(time.Now().Add(time.Second)); err != nil {
		return
	}
	err = wss.ws.WriteMessage(websocket.TextMessage, []byte(message))
	return
}

func (wss *WSSClient) Recv() (message string, err error) {
	if err = wss.ws.SetReadDeadline(time.Now().Add(time.Second)); err != nil {
		return "", err
	}
	_, byteresponse, readerr := wss.ws.ReadMessage()
	if readerr != nil {
		return "", readerr
	}
	return string(byteresponse), nil
}

func (wss *WSSClient) SendRecv(message string) (response string, err error) {
	if err = wss.Send(message); err != nil {
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

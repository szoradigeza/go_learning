package ws

import (
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn               *websocket.Conn
	Pool               *Pool
	FormPoolController *FormPoolController
	Id                 int
}

type Message struct {
	Action string `json:"action"`
	Status string `json:"status"`
	Value  string `json:"value"`
	Id     int    `json:"id"`
}

func getFormData() []FormMock {
	mockData, _ := os.ReadFile("form.json")
	var data []FormMock
	err := json.Unmarshal(mockData, &data)
	if err != nil {
		log.Fatal(err)
	}

	return data
}

type FormMock struct {
	Status string `json:"status"`
	Value  string `json:"value"`
	Id     int    `json:"id"`
}

func (c *Client) Read() {
	defer func() {
		log.Println("defer", c)
		if c.Pool != nil {
			c.Pool.Unregister <- c
		}
		c.Conn.Close()
	}()

	for {
		var message Message

		err := c.Conn.ReadJSON(&message)
		if err != nil {
			log.Println(err)
			return
		}

		log.Printf("Message Received: %+v\n", message)

		if message.Action == "getForm" {
			data := getFormData()[message.Id]
			log.Println("send data ", data)
			c.Conn.WriteJSON(data)
			newClientConnection := ClientConnection{id: message.Id, client: c}
			c.FormPoolController.Register <- &newClientConnection

		} else {
			log.Printf("Broadcast")
			log.Printf(strconv.Itoa(len(c.Pool.Clients)))
			c.Pool.Broadcast <- message
		}
	}
}

package ws

import (
	"fmt"
	"log"
)

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan Message
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
	}
}

func (p *Pool) Start() {
	for {
		select {
		case client := <-p.Register:
			p.Clients[client] = true
			fmt.Println("Size of Connection Pool: ", len(p.Clients))
		case client := <-p.Unregister:
			delete(p.Clients, client)
			if len(client.Pool.Clients) == 0 {
				client.FormPoolController.RemovePool(client.Id)
				log.Println(len(client.Pool.Clients))
			}
		case message := <-p.Broadcast:
			fmt.Println("Size of Connection Pool: ", len(p.Clients))
			for client := range p.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
}

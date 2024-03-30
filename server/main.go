package main

import (
	//"fmt"
	"encoding/json"
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Server struct {
	conns map[*websocket.Conn]bool
}

func newServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) handleWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	s.conns[conn] = true

	if err != nil {
		log.Println(err)
		return
	}

	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}

		log.Printf("Received: %s\n", msg)

		data := BroadcastMessage{}
		json.Unmarshal([]byte(msg), &data)

		broadcastMessage := BroadcastMessage{}
		parseError := json.Unmarshal([]byte(msg), &broadcastMessage)

		if parseError != nil {
			log.Println(parseError)
		}

		log.Println(broadcastMessage)

		s.broadcast(broadcastMessage, conn)

	}
}

type BroadcastMessage struct {
	User  string `json:"user"`
	Value string `json:"value"`
}

func (s *Server) broadcast(broadcastMessage BroadcastMessage, exception *websocket.Conn) {
	byteSlice, err := json.Marshal(broadcastMessage)
	if err != nil {
		log.Println("Error:", err)
		return
	}

	for conn := range s.conns {
		if conn == exception {
			continue
		}
		conn.WriteMessage(websocket.TextMessage, []byte(byteSlice))
	}
}

func main() {
	flag.Parse()
	server := newServer()
	http.HandleFunc("/ws", server.handleWs)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

package main

import (
	//"fmt"

	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Server struct {
	conns map[*websocket.Conn]string
}

func newServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]string),
	}
}

func generateClientID() string {
	uuid := uuid.New()
	return uuid.String()
}

func (s *Server) handleClose(connToRemove *websocket.Conn) {
	connToRemove.Close()
	delete(s.conns, connToRemove)
	fmt.Println(s.conns)
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

func (s *Server) handleWs(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r)
	conn, err := upgrader.Upgrade(w, r, nil)

	s.conns[conn] = generateClientID()

	if err != nil {
		log.Println(err)
		return
	}

	defer s.handleClose(conn)

	for {
		broadcastMessage := BroadcastMessage{}
		err := conn.ReadJSON(&broadcastMessage)
		if err != nil {
			log.Println(err)
			break
		}

		log.Printf("Received: %s\n", broadcastMessage)

		s.broadcast(broadcastMessage, conn)

	}
}

func handleGetForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	data := getFormData()[0]

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	w.Write(jsonBytes)
}

type BroadcastMessage struct {
	User  string `json:"user"`
	Value string `json:"value"`
}

func (s *Server) broadcast(broadcastMessage BroadcastMessage, exception *websocket.Conn) {
	for conn := range s.conns {
		if conn == exception {
			continue
		}

		conn.WriteJSON(broadcastMessage)
	}
}

type FormMock struct {
	Status string `json:"status"`
	Value  string `json:"value"`
}

func main() {
	flag.Parse()
	server := newServer()
	http.HandleFunc("/ws", server.handleWs)
	http.HandleFunc("GET /form/{id}", handleGetForm)

	log.Fatal(http.ListenAndServe(*addr, nil))
}

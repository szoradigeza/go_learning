package main

import (
	//"fmt"

	"flag"
	"fmt"
	"log"
	"net/http"
	ws "server/pkg"

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

func (s *Server) handleWs(fpController *ws.FormPoolController, w http.ResponseWriter, r *http.Request) {
	fmt.Println("handleWs")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	client := ws.Client{
		Conn:               conn,
		FormPoolController: fpController,
	}

	client.Read()
}

func main() {
	flag.Parse()
	server := newServer()
	fpController := ws.NewFormPoolController()
	go fpController.Start()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		server.handleWs(fpController, w, r)
	})

	log.Fatal(http.ListenAndServe(*addr, nil))
}

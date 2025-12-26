package ws

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Server struct {
	Addr string
}

func NewServer(addr string) *Server {
	return &Server{Addr: addr}
}

func (s *Server) Start() error {
	http.HandleFunc("/ws", HandleWebSocket)

	log.Println("[WS] WebSocket iniciado em", s.Addr)
	return http.ListenAndServe(s.Addr, nil)
}

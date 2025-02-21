package main

import (
	"fmt"
	"net"
	"sync"
)

type Client struct {
	conn net.Conn
	send chan string
}

type Server struct {
	clients    map[*Client]bool
	broadcast  chan string
	register   chan *Client
	unregister chan *Client
	mutex      sync.Mutex
}

func NewServer() *Server {
	return &Server{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan string),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (s *Server) Run() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server")
	}
	defer listener.Close()
	fmt.Println("Server is listening of port 8080")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection")
			continue
		}
		client := &Client{conn: conn, send: make(chan string)}
		s.register <- client
		go s.handleClient(client)
	}
}

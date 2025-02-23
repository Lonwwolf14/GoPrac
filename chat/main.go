package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "sync"
    "time"

    "github.com/gorilla/websocket"
)

type Message struct {
    Content   string `json:"content"`
    Sender    string `json:"sender"`
    Timestamp string `json:"timestamp"`
}

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

var clients = make(map[*websocket.Conn]string) // conn -> clientId
var broadcast = make(chan Message)
var mutex = &sync.Mutex{}

func enableCORS(w http.ResponseWriter) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "*")
}

func handler(w http.ResponseWriter, r *http.Request) {
    enableCORS(w)

    if r.Method == "OPTIONS" {
        w.WriteHeader(http.StatusOK)
        return
    }

    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        return
    }
    defer conn.Close()

    mutex.Lock()
    clients[conn] = ""
    mutex.Unlock()

    // Send welcome message
    welcome := Message{
        Content:   "Welcome to the chat!",
        Sender:    "System",
        Timestamp: time.Now().Format("15:04:05"),
    }
    welcomeJSON, _ := json.Marshal(welcome)
    conn.WriteMessage(websocket.TextMessage, welcomeJSON)

    for {
        _, messageBytes, err := conn.ReadMessage()
        if err != nil {
            mutex.Lock()
            delete(clients, conn)
            mutex.Unlock()
            break
        }

        var message Message
        if err := json.Unmarshal(messageBytes, &message); err != nil {
            log.Printf("Error parsing message: %v", err)
            continue
        }

        // Store client ID when first message is received
        if clients[conn] == "" {
            mutex.Lock()
            clients[conn] = message.Sender
            mutex.Unlock()
        }

        broadcast <- message
    }
}

func handleMessages() {
    for {
        message := <-broadcast
        messageJSON, err := json.Marshal(message)
        if err != nil {
            log.Printf("Error marshaling message: %v", err)
            continue
        }

        mutex.Lock()
        for client := range clients {
            // Don't send message back to sender
            if clients[client] != message.Sender {
                err := client.WriteMessage(websocket.TextMessage, messageJSON)
                if err != nil {
                    client.Close()
                    delete(clients, client)
                }
            }
        }
        mutex.Unlock()
    }
}

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        enableCORS(w)
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        http.ServeFile(w, r, "index.html")
    })

    http.HandleFunc("/chat", handler)
    go handleMessages()

    port := ":80"
    fmt.Printf("WebSocket server started on %s\n", port)
    if err := http.ListenAndServe(port, nil); err != nil {
        log.Printf("Error starting server: %v\n", err)
    }
}
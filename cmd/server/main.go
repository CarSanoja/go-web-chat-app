package main

import (
	"log"
	"net/http"
	"sync"

	"go-web-chat-app/config"
	"go-web-chat-app/models"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]string)
var broadcast = make(chan models.Message)
var mutex = &sync.Mutex{}
var maxConnections = 100
var currentConnections = 0

func main() {
	config.LoadConfig()

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	router.GET("/ws", handleConnections)

	go handleMessages()

	log.Printf("Server started at %s:%s", config.GetConfig().ServerAddress, config.GetConfig().ServerPort)
	if err := router.Run(config.GetConfig().ServerAddress + ":" + config.GetConfig().ServerPort); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

func handleConnections(c *gin.Context) {
	mutex.Lock()
	if currentConnections >= maxConnections {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many connections"})
		mutex.Unlock()
		return
	}
	currentConnections++
	mutex.Unlock()

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatalf("Failed to upgrade to WebSocket: %v", err)
	}
	defer ws.Close()

	var username string
	for {
		var msg models.Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error reading JSON: %v", err)
			deleteClient(ws)
			break
		}

		if msg.Type == "username" {
			username = msg.Data
			clients[ws] = username
			continue
		}

		msg.Username = username
		broadcast <- msg
	}

	deleteClient(ws)
}

func deleteClient(ws *websocket.Conn) {
	mutex.Lock()
	defer mutex.Unlock()
	currentConnections--
	delete(clients, ws)
}

func handleMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("Failed to write message to client: %v", err)
				deleteClient(client)
			}
		}
	}
}

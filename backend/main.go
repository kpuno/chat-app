package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Definition of an Upgrader 
// This requires a Read and Write buffer size
var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,

	// Need to check origin of the connection
	// this allows to make request from React development server to here
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Define a reader which will listen for new messages being sent to
// the WebSocket endpoint
func reader(conn *websocket.Conn) {
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		// print out that message for clarity
		fmt.Println(string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}

// define WebSocket endpoint
func serverWs(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Host)

	// upgrade this connection to a WebSocket connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	// listen indefinitely for new messages coming through on the 
	// WebSocket connection
	reader(ws)
}

func setupRoutes() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Simple Server")
	})

	// map our '/ws' endpoint to the 'serverWs' function
	http.HandleFunc("/ws", serverWs)
}


func main() {
	fmt.Println("Chat App v0.01")
	setupRoutes()
  http.ListenAndServe(":8080", nil)
}

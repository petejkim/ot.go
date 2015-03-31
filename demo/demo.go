package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/ws", serveWs)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("public")))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Listening on port %s\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), r)
	if err != nil {
		log.Fatal("Error: ", err)
	}
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		}
		return
	}

	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			return
		}
		if err = conn.WriteMessage(websocket.TextMessage, p); err != nil {
			return
		}
	}
}

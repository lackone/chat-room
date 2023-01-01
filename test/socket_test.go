package test

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"testing"
)

var upgrader = websocket.Upgrader{}

func ws(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func TestWebSocket(t *testing.T) {
	http.HandleFunc("/ws", ws)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

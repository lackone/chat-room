package main

import "github.com/lackone/chat-room/router"

func main() {
	r := router.Router()
	r.Run(":8080")
}

package defines

import "time"

var EmailAddr = "smtp.126.com"
var SendEmail = "lackone@126.com"
var SendEmailPassword = ""

var CodePrefix = "CODE_"
var CodeExpire = time.Minute * 5

// 定义websocket消息结构体
type WebsocketMessage struct {
	From    string `json:"from"`
	To      string `json:"to"`
	RoomId  string `json:"room_id"`
	Message string `json:"message"`
}

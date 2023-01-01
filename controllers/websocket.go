package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/lackone/chat-room/defines"
	"github.com/lackone/chat-room/helpers"
	"github.com/lackone/chat-room/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"
)

var upgrader = websocket.Upgrader{}
var conns = make(map[string]*websocket.Conn)

func WsMessage(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	claims := c.MustGet("jwt_claims")
	jwtClaims := claims.(*helpers.JwtClaims)

	conns[jwtClaims.UserIdentity] = conn

	for {
		wm := defines.WebsocketMessage{}
		err = conn.ReadJSON(&wm)
		if err != nil {
			log.Println(err)
			continue
		}

		//判断用户是否属于当前消息的房间
		userRoomModel := models.UserRoom{}
		_, err = userRoomModel.GetUserRoom(jwtClaims.UserIdentity, wm.RoomId)
		if err != nil {
			log.Println(err)
			continue
		}

		//保存消息
		m := models.Message{
			Id:           primitive.NewObjectID().Hex(),
			UserIdentity: jwtClaims.UserIdentity,
			RoomIdentity: wm.RoomId,
			Data:         wm.Message,
			CreatedAt:    time.Now().Unix(),
			UpdatedAt:    time.Now().Unix(),
		}

		messageModel := models.Message{}
		err = messageModel.InsertOneMessage(&m)
		if err != nil {
			log.Println(err)
			continue
		}

		//获取指定房间在线的用户
		rooms, err := userRoomModel.GetUsersByRooms(wm.RoomId)
		if err != nil {
			log.Println(err)
			continue
		}

		for _, v := range rooms {
			if w, ok := conns[v.UserIdentity]; ok {
				err = w.WriteMessage(websocket.TextMessage, []byte(wm.Message))
				if err != nil {
					log.Println(err)
					continue
				}
			}
		}
	}
}

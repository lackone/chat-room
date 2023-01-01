package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/lackone/chat-room/helpers"
	"github.com/lackone/chat-room/models"
	"net/http"
	"strconv"
)

// 消息列表
func MessageList(c *gin.Context) {
	roomId := c.Query("room_id")

	if roomId == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "房间ID不能为空",
		})
		return
	}

	claims := c.MustGet("jwt_claims")
	jwtClaims := claims.(*helpers.JwtClaims)

	//判断用户是否属于当前消息的房间
	userRoomModel := models.UserRoom{}
	_, err := userRoomModel.GetUserRoom(jwtClaims.UserIdentity, roomId)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "你不在房间内",
		})
		return
	}

	page, _ := strconv.ParseInt(c.Query("page"), 10, 32)
	size, _ := strconv.ParseInt(c.Query("size"), 10, 32)
	skip := (page - 1) * size

	messageModel := models.Message{}

	messages, err := messageModel.GetListByRoomId(roomId, &skip, &size)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "成功",
		"data": gin.H{
			"messages": messages,
		},
	})
}

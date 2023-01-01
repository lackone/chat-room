package router

import (
	"github.com/gin-gonic/gin"
	"github.com/lackone/chat-room/controllers"
	"github.com/lackone/chat-room/middlewares"
)

func Router() *gin.Engine {
	r := gin.Default()

	//登录
	r.POST("/login", controllers.UserLogin)
	//发送验证码
	r.POST("/send_code", controllers.SendCode)
	//注册
	r.POST("/register", controllers.Register)

	auth := r.Group("/user", middlewares.Auth())

	//用户查询
	auth.POST("/query", controllers.UserQuery)

	//用户详情
	auth.GET("/detail", controllers.UserDetail)

	//websocket消息
	auth.GET("/ws/message", controllers.WsMessage)

	//消息列表
	auth.GET("/message/list", controllers.MessageList)

	//添加好友
	auth.POST("/add_friend", controllers.UserAddFriend)

	//删除好友
	auth.POST("/del_friend", controllers.UserDelFriend)

	return r
}

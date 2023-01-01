package controllers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lackone/chat-room/defines"
	"github.com/lackone/chat-room/helpers"
	"github.com/lackone/chat-room/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

// 用户登录
func UserLogin(c *gin.Context) {
	account := c.PostForm("account")
	password := c.PostForm("password")
	if account == "" || password == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "账号或密码不能为空",
		})
		return
	}
	userModel := models.User{}
	user, err := userModel.GetByAccount(account)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "未找到该账号",
		})
		return
	}
	if user.Password != helpers.GetMd5(password) {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "密码错误",
		})
		return
	}
	jwt, err := helpers.MakeJwt(user.Id, user.Account)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "token生成错误",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "成功",
		"data": gin.H{
			"account": user.Account,
			"token":   jwt,
		},
	})
}

// 用户详情
func UserDetail(c *gin.Context) {
	claims := c.MustGet("jwt_claims")
	jwtClaims := claims.(*helpers.JwtClaims)
	userModel := models.User{}
	fmt.Println(jwtClaims.UserIdentity)
	user, err := userModel.GetById(jwtClaims.UserIdentity)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "获取用户失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "成功",
		"data": gin.H{
			"account": user,
		},
	})
}

// 发送验证码
func SendCode(c *gin.Context) {
	email := c.PostForm("email")
	if email == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "email不能为空",
		})
		return
	}
	userModel := models.User{}
	cnt, err := userModel.GetCountByEmail(email)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}
	if cnt > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "邮箱已被注册",
		})
		return
	}
	code := helpers.GetCode()
	err = helpers.EmailSendCode(email, code)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	_, err = models.Redis.Set(context.Background(), defines.CodePrefix+email, code, defines.CodeExpire).Result()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "发送成功",
	})
}

// 注册
func Register(c *gin.Context) {
	email := c.PostForm("email")
	code := c.PostForm("code")
	account := c.PostForm("account")
	password := c.PostForm("password")

	if email == "" || code == "" || account == "" || password == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "参数不完整",
		})
		return
	}

	redisCode, err := models.Redis.Get(context.Background(), defines.CodePrefix+email).Result()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "未找到code",
		})
		return
	}

	if code != redisCode {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "code不正确",
		})
		return
	}

	userModel := models.User{}
	cnt, err := userModel.GetCountByAccount(account)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}
	if cnt > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "账号已存在",
		})
		return
	}

	u := models.User{
		Id:        primitive.NewObjectID().Hex(),
		Account:   account,
		Password:  helpers.GetMd5(password),
		Nickname:  account,
		Sex:       1,
		Email:     email,
		Avatar:    "",
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	err = userModel.InsertOneUser(&u)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "用户创建失败",
		})
		return
	}

	jwt, err := helpers.MakeJwt(u.Id, u.Account)
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
			"account": u.Account,
			"token":   jwt,
		},
	})
}

// 用户查询
func UserQuery(c *gin.Context) {
	account := c.PostForm("account")
	if account == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "账号不能为空",
		})
		return
	}

	userModel := models.User{}
	user, err := userModel.GetByAccount(account)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	claims := c.MustGet("jwt_claims")
	jwtClaims := claims.(*helpers.JwtClaims)

	userRoomModel := models.UserRoom{}
	isFriend, err := userRoomModel.IsFriend(jwtClaims.UserIdentity, user.Id)
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
			"account":   user.Account,
			"nickname":  user.Nickname,
			"avatar":    user.Avatar,
			"sex":       user.Sex,
			"email":     user.Email,
			"is_friend": isFriend,
		},
	})
}

// 添加好友
func UserAddFriend(c *gin.Context) {
	account := c.PostForm("account")
	if account == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "账号不能为空",
		})
		return
	}
	userModel := models.User{}
	user, err := userModel.GetByAccount(account)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}
	claims := c.MustGet("jwt_claims")
	jwtClaims := claims.(*helpers.JwtClaims)
	userRoomModel := models.UserRoom{}
	isFriend, err := userRoomModel.IsFriend(jwtClaims.UserIdentity, user.Id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}
	if isFriend {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "互为好友不可重复添加",
		})
		return
	}
	number := fmt.Sprintf("%d", time.Now().UnixNano())
	r := models.Room{
		Id:           primitive.NewObjectID().Hex(),
		Number:       number,
		Name:         jwtClaims.Account + "创建的" + number + "房间",
		Info:         jwtClaims.Account + "创建的" + number + "房间",
		UserIdentity: jwtClaims.UserIdentity,
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
	}
	roomModel := models.Room{}
	err = roomModel.InsertOneRoom(&r)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}
	ur1 := models.UserRoom{
		Id:           primitive.NewObjectID().Hex(),
		UserIdentity: jwtClaims.UserIdentity,
		RoomIdentity: r.Id,
		RoomType:     2,
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
	}
	err = userRoomModel.InsertOneUserRoom(&ur1)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}
	ur2 := models.UserRoom{
		Id:           primitive.NewObjectID().Hex(),
		UserIdentity: user.Id,
		RoomIdentity: r.Id,
		RoomType:     2,
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
	}
	err = userRoomModel.InsertOneUserRoom(&ur2)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "添加成功",
	})
}

// 删除好友
func UserDelFriend(c *gin.Context) {
	account := c.PostForm("account")
	if account == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "账号不能为空",
		})
		return
	}
	userModel := models.User{}
	user, err := userModel.GetByAccount(account)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}
	claims := c.MustGet("jwt_claims")
	jwtClaims := claims.(*helpers.JwtClaims)
	userRoomModel := models.UserRoom{}
	isFriend, err := userRoomModel.IsFriend(jwtClaims.UserIdentity, user.Id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}
	if !isFriend {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "互相不是好友，不能删除",
		})
		return
	}
	room, err := userRoomModel.GetSingleRoom(jwtClaims.UserIdentity, user.Id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}
	err = userRoomModel.DelByIdAndUser(room, jwtClaims.UserIdentity, user.Id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	roomModel := models.Room{}
	err = roomModel.DelOneRoom(room)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}

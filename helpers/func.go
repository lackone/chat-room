package helpers

import (
	"crypto/md5"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jordan-wright/email"
	"github.com/lackone/chat-room/defines"
	"math/rand"
	"net/smtp"
	"time"
)

// 获取md5
func GetMd5(str string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

// jwt载荷
type JwtClaims struct {
	UserIdentity string `json:"user_identity"`
	Account      string `json:"account"`
	jwt.RegisteredClaims
}

// jwt密钥
var jwtKey = []byte("123456")

// 生成token
func MakeJwt(userIdentity, account string) (string, error) {
	claim := JwtClaims{
		UserIdentity: userIdentity,
		Account:      account,
		RegisteredClaims: jwt.RegisteredClaims{
			//过期时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(3 * time.Hour)),
			//签发时间
			IssuedAt: jwt.NewNumericDate(time.Now()),
			//生效时间
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString(jwtKey)
	return tokenString, err
}

// 解析token
func ParseJwt(tokenString string) (*JwtClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("token valid false")
	}
	if claims, ok := token.Claims.(*JwtClaims); ok {
		return claims, nil
	}
	return nil, errors.New("token error")
}

// 邮箱发送验证码
func EmailSendCode(toEmail, code string) error {
	e := email.NewEmail()
	e.From = "Get <" + defines.SendEmail + ">"
	e.To = []string{toEmail}
	e.Subject = "验证码已发送，请查收"
	e.HTML = []byte("您的验证码为：<b>" + code + "</b")
	return e.SendWithTLS(defines.EmailAddr+":465",
		smtp.PlainAuth("", defines.SendEmail, defines.SendEmailPassword, defines.EmailAddr),
		&tls.Config{InsecureSkipVerify: true, ServerName: defines.EmailAddr},
	)
}

// 获取随机码
func GetCode() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

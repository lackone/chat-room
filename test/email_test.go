package test

import (
	"github.com/lackone/chat-room/helpers"
	"testing"
)

func TestSendEmail(t *testing.T) {
	err := helpers.EmailSendCode("805899763@qq.com", helpers.GetCode())
	t.Log(err)
}

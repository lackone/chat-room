package test

import (
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	t.Log(time.Now().Add(time.Hour * 3).Format("2006-01-02 15:04:05"))
}

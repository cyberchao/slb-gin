package main

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
)

type Person struct {
	Username string `form:"username"`
	Password string `form:"password"`
	Iotp     string `form:"iotp"`
	Sms      string `form:"sms"`
}

func main() {
	route := gin.Default()
	route.POST("/api/v1.2/gen_token_by_param", token)
	route.Run(":8085")
}

func gen_token(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ="
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func token(c *gin.Context) {
	person1 := Person{
		Username: "jiazhichao621",
		Password: "Qwer4321!",
		Iotp:     "777888",
		Sms:      "",
	}
	var person Person

	if c.ShouldBind(&person) == nil {
		if person == person1 {
			c.JSON(http.StatusOK, gin.H{"returnCode": "0", "token": gen_token(48)})
		} else {
			c.JSON(http.StatusOK, gin.H{"returnCode": "1", "token": ""})
		}
	} else {
		c.JSON(http.StatusOK, gin.H{"returnCode": "2", "token": ""})
	}

}

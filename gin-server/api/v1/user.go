package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	type RegisterAndLoginStruct struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Iotp     string `json:"iotp"`
		Sms      string `json:"sms"`
	}
	var L RegisterAndLoginStruct
	_ = c.ShouldBindJSON(&L)
	fmt.Println("login page")

}

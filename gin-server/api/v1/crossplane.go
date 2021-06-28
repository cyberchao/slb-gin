package v1

import (
	"github.com/aluttik/go-crossplane"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"slb-admin/global/response"
	"strings"
	"time"
)

func CrossParse(c *gin.Context) {
	requestData, _:=c.GetRawData()
	ngxConf := string(requestData)
	ngxConf = strings.ReplaceAll(strings.Trim(ngxConf,"\""),"\\n","\n")


	ioutil.WriteFile("tmp/data.txt",[]byte(ngxConf), 0644)

	c1 := make(chan crossplane.Payload, 1)
	go func() {
		payload, _ := crossplane.Parse("tmp/data.txt", &crossplane.ParseOptions{})
		c1 <- *payload
	}()

	select {
	case res := <-c1:
		if res.Status == "ok" {
			response.OkWithData(res.Config[0].Parsed[0], c)
		}else{
			response.FailWithMessage("解析失败",c)
		}
	case <-time.After(1 * time.Second):
		response.FailWithMessage("解析超时",c)
	}
}
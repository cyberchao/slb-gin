package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"slb-admin/model/request"
)

func LoginSSO(L request.LoginStruct) (err error, ssores map[string]interface{}) {

	data := make(url.Values)
	data["username"] = []string{L.Username}
	data["password"] = []string{L.Password}
	data["iotp"] = []string{L.Iotp}
	data["sms"] = []string{L.Sms}
	var res map[string]interface{}
	response, err := http.PostForm("http://127.0.0.1:8085/api/v1.2/gen_token_by_param", data)
	if err != nil {
		fmt.Println(err)
		return err, res
	}
	defer response.Body.Close()
	json.NewDecoder(response.Body).Decode(&res)
	return nil, res
	//if res["returnCode"] == "1" {
	//	return err, ""
	//} else {
	//	return err, res["token"].(string)
	//}
}

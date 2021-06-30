package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"slb-admin/global"
	"slb-admin/global/response"
	"slb-admin/model"
	resp "slb-admin/model/response"
	"strings"
	"time"

	"github.com/aluttik/go-crossplane"
	"github.com/apenella/go-ansible/pkg/adhoc"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Tags vhost
// @Summary 新建vhost
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"新建vhost成功"}"
// @Router /role/deleteRole [post]
func CreateServer(c *gin.Context) {
	type requestData struct {
		Env         string `json:"env"`
		Cluster     string `json:"cluster"`
		NgxConf     string `json:"code"`
		Description string `json:"description"'`
		Version     int    `json:"version"`
		FileName    string `json:"filename"`
	}

	var postData requestData
	if err := c.ShouldBindJSON(&postData); err != nil {
		response.FailWithMessage("请求数据异常", c)
	}
	//requestData, _:=c.GetRawData()
	ngxConf := postData.NgxConf
	ngxConf = strings.ReplaceAll(strings.Trim(ngxConf, "\""), "\\n", "\n")
	http_ngxConf := "http {\n" + ngxConf + "\n}"

	ioutil.WriteFile("tmp/data.txt", []byte(http_ngxConf), 0644)

	c1 := make(chan crossplane.Payload, 1)
	go func() {
		payload, _ := crossplane.Parse("tmp/data.txt", &crossplane.ParseOptions{})
		c1 <- *payload //通过管道设置超时时间为1s
	}()

	select {
	case res := <-c1:
		if res.Status == "ok" {
			jsonString := *(res.Config[0].Parsed[0].Block)
			b, _ := json.Marshal(jsonString[0])

			var ngxdoc interface{}

			json.Unmarshal(b, &ngxdoc)
			serverArgs := gjson.Get(string(b), "block.#(directive==\"server_name\").args").Array()

			portArgs := gjson.Get(string(b), "block.#(directive==\"listen\").args").Array()
			fmt.Println(serverArgs, portArgs)

			collection := global.Mogo.Database("slb").Collection("vhost")
			//检查域名和端口是否已存在 db.vhost.find({"ngx.block": { $elemMatch: {args:"yqb.com"}}})
			for _, val := range serverArgs {
				serverName := val.String()
				filter := bson.M{
					"ngx.block": bson.M{
						"$elemMatch": bson.M{"args": serverName}},
					"cluster": postData.Cluster, "env": postData.Env}
				filterCursor, err := collection.Find(context.TODO(), filter)
				var result []bson.M
				if err = filterCursor.All(context.TODO(), &result); err != nil {
					log.Fatal(err)
				}
				if result != nil {
					filter = bson.M{
						"ngx.block": bson.M{
							"$elemMatch": bson.M{"args": portArgs[0]}},
						"cluster": postData.Cluster, "env": postData.Env}
					var result []bson.M
					if err = filterCursor.All(context.TODO(), &result); err != nil {
						log.Fatal(err)
					}
					if result != nil {
						response.FailWithMessage("域名已存在:"+serverName, c)
						return
					}
				}
			}

			basepath := global.CONFIG.System.Path
			filepath := fmt.Sprintf("%s/%s/%s/vhost", basepath, postData.Env, postData.Cluster)
			doc := model.VhostDoc{
				Env:         postData.Env,
				Cluster:     postData.Cluster,
				Ngx:         ngxdoc,
				Src:         postData.NgxConf,
				Description: postData.Description,
				Status:      false,
				Version:     postData.Version,
				Time:        time.Now(),
				FilePath:    filepath + "/" + postData.FileName}
			collection.InsertOne(context.TODO(), doc)
			err := os.MkdirAll(filepath, 0755)
			if err != nil {
				log.Fatal(err)
			}
			//写入conf文件，相同的域名写入同一个文件
			f, err := os.OpenFile(filepath+"/"+postData.FileName,
				os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				response.FailWithMessage(err.Error(), c)
				return
			}
			defer f.Close()
			if _, err := f.WriteString(ngxConf + "\n"); err != nil {
				response.FailWithMessage(err.Error(), c)
				return
			}
			//err = ioutil.WriteFile(filepath+"/"+postData.FileName, []byte(ngxConf), 0644)
			//if err != nil {
			//	log.Fatal(err)
			//}

			response.OkWithMessage("success", c)
		} else {
			response.FailWithMessage("格式错误-0", c)
		}
	case <-time.After(1 * time.Second):
		response.FailWithMessage("格式错误-1", c)
	}
}
func UpdateServer(c *gin.Context) {
	type requestData struct {
		Newcode     string `json:"newcode"`
		Id          string `json:"id"`
		Env         string `json:"env"`
		Cluster     string `json:"cluster"`
		Description string `json:"description"`
		Filename    string `json:"filename"`
		Version     int    `json:"version"`
	}

	var postData requestData
	if err := c.ShouldBindJSON(&postData); err != nil {
		response.FailWithMessage("请求数据异常", c)
	}

	ngxConf := strings.ReplaceAll(strings.Trim(postData.Newcode, "\""), "\\n", "\n")
	http_ngxConf := "http {\n" + ngxConf + "\n}"

	ioutil.WriteFile("tmp/data.txt", []byte(http_ngxConf), 0644)

	c1 := make(chan crossplane.Payload, 1)
	go func() {
		payload, _ := crossplane.Parse("tmp/data.txt", &crossplane.ParseOptions{})
		c1 <- *payload //通过管道设置超时时间为1s
	}()

	select {
	case res := <-c1:
		if res.Status == "ok" {
			jsonString := *(res.Config[0].Parsed[0].Block)
			b, _ := json.Marshal(jsonString[0])

			var ngxdoc interface{}

			json.Unmarshal(b, &ngxdoc)
			serverArgs := gjson.Get(string(b), "block.#(directive==\"server_name\").args").Array()

			portArgs := gjson.Get(string(b), "block.#(directive==\"listen\").args").Array()
			fmt.Println(serverArgs, portArgs)

			collection := global.Mogo.Database("slb").Collection("vhost")
			//检查域名和端口是否已存在 db.vhost.find({"ngx.block": { $elemMatch: {args:"yqb.com"}}})
			for _, val := range serverArgs {
				serverName := val.String()
				filter := bson.M{
					"ngx.block": bson.M{
						"$elemMatch": bson.M{"args": serverName}},
					"cluster": postData.Cluster, "env": postData.Env}
				filterCursor, err := collection.Find(context.TODO(), filter)
				var result []bson.M
				if err = filterCursor.All(context.TODO(), &result); err != nil {
					log.Fatal(err)
				}
				if result != nil {
					filter = bson.M{
						"ngx.block": bson.M{
							"$elemMatch": bson.M{"args": portArgs[0]}},
						"cluster": postData.Cluster, "env": postData.Env}
					var result []bson.M
					if err = filterCursor.All(context.TODO(), &result); err != nil {
						log.Fatal(err)
					}
					if result != nil {
						response.FailWithMessage("域名已存在:"+serverName, c)
						return
					}
				}
			}

			doc := model.VhostDoc{
				Env:         postData.Env,
				Cluster:     postData.Cluster,
				Ngx:         ngxdoc,
				Src:         postData.Newcode,
				Description: postData.Description,
				Status:      false,
				Version:     postData.Version,
				Time:        time.Now(),
				FilePath:    postData.Filename}

			id, _ := primitive.ObjectIDFromHex(postData.Id)
			collection.ReplaceOne(
				context.TODO(),
				bson.M{"_id": id},
				doc,
			)

			//覆盖写入conf文件
			f, err := os.OpenFile(postData.Filename,
				os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
			if err != nil {
				response.FailWithMessage(err.Error(), c)
				return
			}
			defer f.Close()
			if _, err := f.WriteString(ngxConf + "\n"); err != nil {
				response.FailWithMessage(err.Error(), c)
				return
			}

			response.OkWithMessage("success", c)
		} else {
			response.FailWithMessage("格式错误-0", c)
		}
	case <-time.After(1 * time.Second):
		response.FailWithMessage("格式错误-1", c)
	}
}

func DeleteServer(c *gin.Context) {
	type requestData struct {
		Id string `json:"id"`

		FilePath string `json:"filepath"`
	}
	var postData requestData
	if err := c.ShouldBindJSON(&postData); err != nil {
		response.FailWithMessage("请求数据异常", c)
	}
	collection := global.Mogo.Database("slb").Collection("vhost")

	idPrimitive, _ := primitive.ObjectIDFromHex(postData.Id)
	res, err := collection.DeleteOne(context.TODO(), bson.M{"_id": idPrimitive})
	if err != nil {
		response.FailWithMessage("删除失败", c)
		return
	}
	// Check if the response is 'nil'
	if res.DeletedCount == 0 {
		response.FailWithMessage("server不存在", c)

	} else {
		err = os.Remove(postData.FilePath)
		if err != nil {
			response.FailWithMessage(err.Error(), c)

		}
		response.OkWithMessage("删除成功", c)
	}
}

func GetServerList(c *gin.Context) {
	type requestData struct {
		Page        int      `json:"page"`
		PageSize    int      `json:"pageSize"`
		Env         []string `json:"env"`
		Cluster     []string `json:"cluster"`
		Server_name string   `json:"server_name"`
	}

	var postData requestData
	if err := c.ShouldBindJSON(&postData); err != nil {
		response.FailWithMessage("请求数据异常", c)
	}
	collection := global.Mogo.Database("slb").Collection("vhost")

	//设置分页条件
	findOptions := options.Find()
	findOptions.SetLimit(int64(postData.PageSize))
	findOptions.SetSkip(int64(postData.PageSize * (postData.Page - 1)))

	//条件查询 https://stackoverflow.com/questions/55306617/how-to-add-values-to-an-bson-d-object
	var filter = bson.D{}
	if len(postData.Env) != 0 {
		filter = append(filter, bson.E{"env", bson.M{"$in": postData.Env}})
	}
	if len(postData.Cluster) != 0 {
		filter = append(filter, bson.E{"cluster", bson.M{"$in": postData.Cluster}})
	}
	if postData.Server_name != "" {
		filter = append(filter, bson.E{"ngx.block", bson.M{
			"$elemMatch": bson.M{"args": postData.Server_name}}})
	}

	count, _ := collection.CountDocuments(context.TODO(), filter)
	filterCursor, _ := collection.Find(context.TODO(), filter, findOptions)

	var result []bson.M
	if err := filterCursor.All(context.TODO(), &result); err != nil {
		log.Fatal(err)
	}
	response.OkWithData(resp.PageResult{
		List:     result,
		Total:    count,
		Page:     postData.Page,
		PageSize: postData.PageSize,
	}, c)
}

func PublishServer(c *gin.Context) {
	type requestData struct {
		Id       string `json:"id"`
		Env      string `json:"env"`
		Cluster  string `json:"cluster"`
		Filepath string `json:"filepath"`
	}

	var postData requestData
	if err := c.ShouldBindJSON(&postData); err != nil {
		response.FailWithMessage("请求数据异常", c)
	}

	db := global.DB
	var hosts []model.Host

	db.Where("env = ? AND cluster = ?", postData.Env, postData.Cluster).Find(&hosts)
	ipstr := ""
	for _, v := range hosts {
		ipstr += v.Ip + ","
	}
	remotePath := global.CONFIG.System.RemotePath
	args := fmt.Sprintf("src=%s dest=%s", postData.Filepath, remotePath)

	ansibleAdhocOptions := &adhoc.AnsibleAdhocOptions{
		Inventory:  ipstr,
		ModuleName: "synchronize",
		Args:       args,
	}

	adhoc := &adhoc.AnsibleAdhocCmd{
		Pattern: "all",
		Options: ansibleAdhocOptions,
	}

	fmt.Println("adhoc:", adhoc.String())
	err := adhoc.Run(context.TODO())
	if err != nil {
		panic(err)
	}
	collection := global.Mogo.Database("slb").Collection("vhost")

	idPrimitive, _ := primitive.ObjectIDFromHex(postData.Id)
	_, err = collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": idPrimitive},
		bson.M{"$set": bson.M{"status": true}},
	)
	if err != nil {
		global.Logger.Errorf("edit vhost status failed, err:[%s] ", err.Error())
		response.FailWithMessage(err.Error(), c)
	}

	response.OkWithMessage("发布成功", c)
}

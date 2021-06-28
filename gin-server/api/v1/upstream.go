package v1

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"slb-admin/global"
	"slb-admin/global/response"
	"slb-admin/model"
	resp "slb-admin/model/response"
	"time"
)

func CreateUpstream(c *gin.Context) {

	type requestData struct {
		Env        string             `json:"env"`
		Cluster    string             `json:"cluster"`
		Name       string             `json:"name"`
		ServerList []model.ServerList `json:"serverList"'`
		Forward    string             `json:"forward"`
	}

	var postData requestData
	if err := c.ShouldBindJSON(&postData); err != nil {
		response.FailWithMessage("请求数据异常", c)
	}

	collection := global.Mogo.Database("slb").Collection("upstream")

	filter := bson.M{
		"name":    postData.Name,
		"cluster": postData.Cluster, "env": postData.Env}
	filterCursor, err := collection.Find(context.TODO(), filter)
	var result []bson.M
	if err = filterCursor.All(context.TODO(), &result); err != nil {
		log.Fatal(err)
	}
	if result != nil {
		response.FailWithMessage("upstream已存在:"+postData.Name, c)
		return
	}
	basepath := global.CONFIG.System.Path
	filepath := fmt.Sprintf("%s/%s/%s/upstream", basepath, postData.Env, postData.Cluster)
	doc := model.UpstreamDoc{
		Env:        postData.Env,
		Cluster:    postData.Cluster,
		Name:       postData.Name,
		ServerList: postData.ServerList,
		Forward:    postData.Forward,
		Version:    1,
		Time:       time.Now(),
		FilePath:   filepath + "/" + postData.Name + ".conf"}
	collection.InsertOne(context.TODO(), doc)
	err = os.MkdirAll(filepath, 0755)
	if err != nil {
		log.Fatal(err)
	}
	var confTmp string
	if postData.Forward == "chash" {
		confTmp = fmt.Sprintf(`upstream %[1]s {
 		server 0.0.0.1;
        balancer_by_lua_block {
            local b = require "ngx.balancer"
            local chash_up = package.loaded.%[1]s_chash_up
            local servers = package.loaded.%[1]s_servers
            local id = chash_up:find(ngx.var.arg_key)
            local server = servers[id]
            assert(b.set_current_peer(server))
        }
}`, postData.Name)
	} else if postData.Forward == "roundbin" {
		confTmp = fmt.Sprintf(`upstream %[1]s {
   	server 0.0.0.1;
    balancer_by_lua_block {
        local b = require "ngx.balancer"
        local rr_up = package.loaded.%[1]s_rr_up
        local server = rr_up:find()
        assert(b.set_current_peer(server))
    }
}`, postData.Name)
	} else {
		response.FailWithMessage("生成配置文件失败", c)
		return
	}

	f, err := os.OpenFile(filepath+"/"+postData.Name+".conf",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	defer f.Close()
	if _, err := f.WriteString(confTmp + "\n"); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	//err = ioutil.WriteFile(filepath+"/"+postData.Name+".conf", []byte(confTmp), 0644)
	//if err != nil {
	//	response.FailWithMessage(err.Error(), c)
	//	return
	//}
	response.OkWithMessage("success", c)
}
func UpdateUpstream(c *gin.Context) {
	response.FailWithMessage("解析失败", c)
}

func DeleteUpstream(c *gin.Context) {
	type requestData struct {
		Id       string `json:"id"`
		FilePath string `json:"filepath"`
	}
	var postData requestData
	if err := c.ShouldBindJSON(&postData); err != nil {
		response.FailWithMessage("请求数据异常", c)
		return
	}
	collection := global.Mogo.Database("slb").Collection("upstream")

	idPrimitive, err := primitive.ObjectIDFromHex(postData.Id)

	res, err := collection.DeleteOne(context.TODO(), bson.M{"_id": idPrimitive})
	if err != nil {
		response.FailWithMessage("删除失败", c)
		return
	}
	// Check if the response is 'nil'
	if res.DeletedCount == 0 {
		response.FailWithMessage("upstream不存在", c)
	} else {
		err = os.Remove(postData.FilePath)
		if err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
		response.OkWithMessage("删除成功", c)
	}
}

func GetUpstreamList(c *gin.Context) {
	type requestData struct {
		Page     int      `json:"page"`
		PageSize int      `json:"pageSize"`
		Env      []string `json:"env"`
		Cluster  []string `json:"cluster"`
		Name     string   `json:"name"`
	}

	var postData requestData
	if err := c.ShouldBindJSON(&postData); err != nil {
		response.FailWithMessage("请求数据异常", c)
	}
	collection := global.Mogo.Database("slb").Collection("upstream")

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
	if postData.Name != "" {
		filter = append(filter, bson.E{"ngx.block", bson.M{
			"$elemMatch": bson.M{"args": postData.Name}}})
	}

	count, err := collection.CountDocuments(context.TODO(), filter)
	filterCursor, err := collection.Find(context.TODO(), filter, findOptions)

	var result []bson.M
	if err = filterCursor.All(context.TODO(), &result); err != nil {
		log.Fatal(err)
	}

	response.OkWithData(resp.PageResult{
		List:     result,
		Total:    count,
		Page:     postData.Page,
		PageSize: postData.PageSize,
	}, c)
}

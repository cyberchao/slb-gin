package v1

import (
	"context"
	"fmt"
	"os"
	"slb-admin/global"
	"slb-admin/global/response"
	"slb-admin/model"
	resp "slb-admin/model/response"
	"time"

	"github.com/apenella/go-ansible/pkg/adhoc"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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
		global.Logger.Errorf("c.ShouldBindJSON failed. err: [%s]", err.Error())
		response.FailWithMessage("请求数据异常", c)
	}

	collection := global.Mogo.Database("slb").Collection("upstream")

	filter := bson.M{
		"name":    postData.Name,
		"cluster": postData.Cluster, "env": postData.Env}
	filterCursor, err := collection.Find(context.TODO(), filter)
	var result []bson.M
	if err = filterCursor.All(context.TODO(), &result); err != nil {
		global.Logger.Errorf("mongo filterCursor. err: [%s]", err.Error())
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
	var res interface{}
	res, err = collection.InsertOne(context.TODO(), doc)
	if err != nil {
		response.FailWithMessage("插入数据库异常", c)
		global.Logger.Errorf("insert to mongo failed. err: [%s]", err.Error())
	}

	global.Logger.Infof("insert to mongo success. result: %s", res)
	err = os.MkdirAll(filepath, 0755)
	if err != nil {
		global.Logger.Errorf("os.MkdirAll. err: [%s]", err.Error())
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
		os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		global.Logger.Errorf("open file failed. err: [%s]", err.Error())
		response.FailWithMessage(err.Error(), c)
		return
	}
	defer f.Close()
	if _, err := f.WriteString(confTmp + "\n"); err != nil {
		global.Logger.Errorf("write upstream conf [%s] to file [%s] failed %s. err: [%s]", confTmp, f, err.Error())
		response.FailWithMessage(err.Error(), c)
		return
	}
	db := global.DB
	var hosts []model.Host

	db.Where("env = ? AND cluster = ?", postData.Env, postData.Cluster).Find(&hosts)
	ipstr := ""
	for _, v := range hosts {
		ipstr += v.Ip + ","
	}
	remotePath := global.CONFIG.System.RemotePath
	args := fmt.Sprintf("src=%s dest=%s", doc.FilePath, remotePath+"/upstream")
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
	err = adhoc.Run(context.TODO())
	if err != nil {
		global.Logger.Errorf("ansible adhoc run error,%s ", err.Error())
	}
	global.Logger.Infof("create upstream success,%s ", doc)
	response.OkWithMessage("success", c)
}
func UpdateUpstream(c *gin.Context) {
	type requestData struct {
		Id            string             `json:"id"`
		NewServerList []model.ServerList `json:"form"`
	}
	var postData requestData
	if err := c.ShouldBindJSON(&postData); err != nil {
		global.Logger.Errorf("c.ShouldBindJSON failed. err: [%s]", err.Error())
		response.FailWithMessage("请求数据异常", c)
		return
	}

	collection := global.Mogo.Database("slb").Collection("upstream")
	id, _ := primitive.ObjectIDFromHex(postData.Id)
	result, err := collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"serverlist": postData.NewServerList}},
	)
	if err != nil {
		global.Logger.Errorf("update upstream failed, err:[%s] ", err.Error())
		response.FailWithMessage(err.Error(), c)
	}
	global.Logger.Infof("update upstream success, err:[%s] ", result)
	response.OkWithMessage("更新成功", c)
}

func DeleteUpstream(c *gin.Context) {
	type requestData struct {
		Id       string `json:"id"`
		FilePath string `json:"filepath"`
	}

	var postData requestData
	if err := c.ShouldBindJSON(&postData); err != nil {
		global.Logger.Errorf("c.ShouldBindJSON failed. err: [%s]", err.Error())
		response.FailWithMessage("请求数据异常", c)
		return
	}
	collection := global.Mogo.Database("slb").Collection("upstream")

	idPrimitive, _ := primitive.ObjectIDFromHex(postData.Id)

	// 根据id获取upstream的env和cluster等信息，生成ansible语句
	result := collection.FindOne(context.Background(), bson.M{"_id": idPrimitive})
	doc1 := model.UpstreamDoc{}
	result.Decode(doc1)
	env := doc1.Env
	cluster := doc1.Cluster
	name := doc1.Name

	res, err := collection.DeleteOne(context.TODO(), bson.M{"_id": idPrimitive})
	if err != nil {
		global.Logger.Errorf("delete upstream failed, err:[%s] ", err.Error())
		response.FailWithMessage("删除失败", c)
		return
	}
	// Check if the response is 'nil'
	if res.DeletedCount == 0 {
		response.FailWithMessage("upstream不存在", c)
	} else {
		err = os.Remove(postData.FilePath)
		if err != nil {
			global.Logger.Errorf("delete upstream file failed, err:[%s] ", err.Error())
			response.FailWithMessage(err.Error(), c)
		}

		db := global.DB
		var hosts []model.Host

		db.Where("env = ? AND cluster = ?", env, cluster).Find(&hosts)
		ipstr := ""
		for _, v := range hosts {
			ipstr += v.Ip + ","
		}
		remotePath := global.CONFIG.System.RemotePath
		args := fmt.Sprintf("dest=%s state=absent", remotePath+"/upstream/"+name+".conf")
		ansibleAdhocOptions := &adhoc.AnsibleAdhocOptions{
			Inventory:  ipstr,
			ModuleName: "ansible.builtin.file",
			Args:       args,
		}

		adhoc := &adhoc.AnsibleAdhocCmd{
			Pattern: "all",
			Options: ansibleAdhocOptions,
		}

		fmt.Println("adhoc:", adhoc.String())
		err = adhoc.Run(context.TODO())
		if err != nil {
			global.Logger.Errorf("ansible adhoc run error,%s ", err.Error())
		}

		global.Logger.Infof("delete upstream success, name: [%s] ", name)
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
		global.Logger.Errorf("c.ShouldBindJSON failed. err: [%s]", err.Error())
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

	count, _ := collection.CountDocuments(context.TODO(), filter)
	filterCursor, _ := collection.Find(context.TODO(), filter, findOptions)

	var result []bson.M
	if err := filterCursor.All(context.TODO(), &result); err != nil {
		global.Logger.Errorf("mongo filterCursor. err: [%s]", err.Error())
	}
	global.Logger.Info("get upstream list success")
	response.OkWithData(resp.PageResult{
		List:     result,
		Total:    count,
		Page:     postData.Page,
		PageSize: postData.PageSize,
	}, c)
}

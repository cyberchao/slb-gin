package v1

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slb-admin/global"
	"slb-admin/model"
	"slb-admin/service"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
)

type Message struct {
	Host    string
	Content string
}
type Server struct {
	Hostname string
	File     string
	Key      string
	Stdout   io.Reader
}

type readJson struct {
	Type      string `json:"type"`
	Key       string `json:"key"`
	Env       string `json:"env"`
	Cluster   string `json:"cluster"`
	AccessLog string `json:"access_log"`
	ErrorLog  string `json:"error_log"`
}

// Execute the remote command
func (server *Server) Execute(output chan Message) {

	user := global.CONFIG.Ssh.User
	port := global.CONFIG.Ssh.Port
	keypath := global.CONFIG.Ssh.KeyPath
	client, _ := service.NewSshClient(
		user,
		server.Hostname,
		port,
		keypath)
	session := client.SshSession()
	TerminalModes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	session.RequestPty("xterm", 80, 40, TerminalModes)
	server.Stdout, _ = session.StdoutPipe()

	go tailOutput(server.Hostname, output, &server.Stdout)
	//session.Start("tail -f " + server.File)

	if err := session.Start("tail -f " + server.File + " " + server.Key); err != nil {
		global.Logger.Errorf(fmt.Sprintf("[%s] failed to execute command: %s", server.Hostname, err))
	}
	session.Wait()
}

// bing the pipe output for formatted output to channel
func tailOutput(host string, output chan Message, input *io.Reader) {
	reader := bufio.NewReader(*input)
	for {
		line, _ := reader.ReadString('\n')
		// global.Logger.Info("read message", line)
		output <- Message{
			Host:    host,
			Content: line,
		}
	}
}

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func RemoteTail(c *gin.Context) {
	//Upgrade get request to webSocket protocol
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		global.Logger.Errorf("error get connection ", err.Error())
	}
	defer ws.Close()

	var readJson readJson
	ws.ReadJSON(&readJson)
	var hosts []model.Host

	// 获取域名对应的nginx主机
	global.DB.Where("cluster = ? AND env = ?", readJson.Cluster, readJson.Env).Find(&hosts)
	var serverList []Server
	for _, v := range hosts {
		var server Server
		server.Hostname = v.Ip
		server.Key = readJson.Key
		if readJson.Type == "access" {
			server.File = readJson.AccessLog
		} else {
			server.File = readJson.ErrorLog
		}
		serverList = append(serverList, server)
	}

	outputs := make(chan Message, 255)
	//等待goroutines执行结束 https://gobyexample.com/waitgroups
	var wg sync.WaitGroup

	for _, server := range serverList {
		wg.Add(1)
		t1, _ := json.Marshal(server)
		global.Logger.Info("start tail", string(t1))
		go func(server Server) {
			defer wg.Done()
			server.Execute(outputs)
		}(server)
	}

	go func() {
		for output := range outputs {
			writedata, _ := json.Marshal(output)
			ws.WriteMessage(1, []byte(writedata))
		}
	}()

	wg.Wait()
}

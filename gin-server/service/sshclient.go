package service

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"slb-admin/global"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

type SshClient struct {
	Config *ssh.ClientConfig
	Server string
}

func NewSshClient(user string, host string, port int, privateKeyPath string) (*SshClient, error) {
	// read private key file
	key, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("Reading private key file failed %v", err)
	}
	// create signer
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}
	// build SSH client config
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			// use OpenSSH's known_hosts file if you care about host validation
			return nil
		},
		Timeout: 3 * time.Second,
	}

	client := &SshClient{
		Config: config,
		Server: fmt.Sprintf("%v:%v", host, port),
	}
	return client, nil
}

func (s *SshClient) RunCommand(cmd string) (string, error) {
	// open connection
	conn, err := ssh.Dial("tcp", s.Server, s.Config)
	if err != nil {
		return "", fmt.Errorf("Dial to %v failed %v", s.Server, err)
	}
	defer conn.Close()

	// open session
	session, err := conn.NewSession()
	if err != nil {
		return "", fmt.Errorf("Create session for %v failed %v", s.Server, err)
	}
	defer session.Close()

	// run command and capture stdout/stderr
	output, err := session.CombinedOutput(cmd)

	return fmt.Sprintf("%s", output), err
}

func SshWoker(hostList []string, cmd string) (map[string]string, error) {

	user := global.CONFIG.Ssh.User
	port := global.CONFIG.Ssh.Port
	keypath := global.CONFIG.Ssh.KeyPath
	res := make(map[string]string)
	for _, ip := range hostList {
		ssh, err := NewSshClient(
			user,
			ip,
			port,
			keypath)

		if err != nil {
			log.Printf("SSH init error %v", err)
		} else {
			output, err := ssh.RunCommand(cmd)
			// 超时错误和其他详细错误区分开
			if err != nil {
				if strings.Contains(err.Error(), "i/o timeout") {
					res[ip] = err.Error()
				} else {
					res[ip] = output
				}
			}
		}
	}
	return res, nil

}

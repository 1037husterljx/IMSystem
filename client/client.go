package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	connect    net.Conn
	Name       string
	ServerIP   string
	ServerPort int
	flag       int
}

func (this *Client) menu() bool {
	var flag int
	fmt.Println("1.公聊")
	fmt.Println("2.私聊")
	fmt.Println("3.更新用户名")
	fmt.Println("0.推出")
	fmt.Scanln(&flag)
	if flag >= 0 && flag <= 3 {
		this.flag = flag
		return true
	} else {
		fmt.Println("请输入合法数字")
		return false
	}
}

var serverIP string
var serverPort int

func NewClient(serverIP string, serverPort int) *Client {
	client := &Client{
		ServerIP:   serverIP,
		ServerPort: serverPort,
		flag:       999,
	}

	connect, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIP, serverPort))

	if err != nil {
		fmt.Println("dial err:", err)
		return nil
	}

	client.connect = connect

	return client
}

func init() {
	flag.StringVar(&serverIP, "ip", "127.0.0.1", "设置服务器ip，默认127.0.0.1")
	flag.IntVar(&serverPort, "port", 8888, "设置端口，默认8888")
}

func (this *Client) Run() {
	for this.flag != 0 {
		for this.menu() != true {
		}
		switch this.flag {
		case 1:
			this.publicChat()
		case 2:
			this.privateChat()
		case 3:
			this.updateUsername()
		case 0:
			fmt.Println("退出")
		}
	}
}

func (this *Client) privateChat() {
	var name string
	var msg string

	this.who()
	fmt.Println("输入用户名,exit退出")
	fmt.Scanln(&name)

	for name != "exit" {
		fmt.Println("输入消息，exit退出")
		fmt.Scanln(&msg)
		for msg != "exit" {
			if len(msg) != 0 {
				msg = "to|" + name + "|" + msg
				_, err := this.connect.Write([]byte(msg))
				if err != nil {
					fmt.Println("connect err:", err)
					break
				}
			}
			msg = ""
			fmt.Println("输入消息，exit退出")
			fmt.Scanln(&msg)
		}

	}
}

func (this *Client) who() {
	msg := "who"

	_, err := this.connect.Write([]byte(msg))
	if err != nil {
		fmt.Println("coonect err:", err)
		return
	}

}

func (this *Client) publicChat() {
	var msg string
	fmt.Println("输入消息，exit退出")

	fmt.Scanln(&msg)

	for msg != "exit" {
		if len(msg) != 0 {
			_, err := this.connect.Write([]byte(msg))
			if err != nil {
				fmt.Println("connect err:", err)
				break
			}
		}
		msg = ""
		fmt.Println("输入消息，exit退出")
		fmt.Scanln(&msg)

	}
}

func (this *Client) printRespose() {
	io.Copy(os.Stdout, this.connect)
}

func (this *Client) updateUsername() bool {
	fmt.Println("请输入用户名")
	fmt.Scanln(&this.Name)

	msg := "rename|" + this.Name

	_, err := this.connect.Write([]byte(msg))
	if err != nil {
		fmt.Println("coonect err:", err)
		return false
	}

	return true
}

func main() {
	flag.Parse()

	client := NewClient(serverIP, serverPort)

	if client == nil {
		fmt.Println("连接失败")
		return
	}

	fmt.Println("连接成功")

	go client.printRespose()

	client.Run()
}

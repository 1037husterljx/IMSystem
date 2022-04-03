package main

import (
	"fmt"
	"net"
)

type Client struct {
	connect    net.Conn
	Name       string
	ServerIP   string
	ServerPort int
}

func NewClient(serverIP string, serverPort int) *Client {
	client := &Client{
		ServerIP:   serverIP,
		ServerPort: serverPort,
	}
	connect, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIP, serverPort))

	if err!=nil {
		fmt.Println("dial err:",err)
		return nil
	}

	client.connect = connect

	return client
}

func main()  {
	client:=NewClient("127.0.0.1",8888)
	if client==nil {
		fmt.Println("连接失败")
		return
	}

	fmt.Println("连接成功")

	select{}
}

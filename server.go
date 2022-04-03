package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip      string
	Port    int
	Lock    sync.RWMutex
	Map     map[string]*User
	Channle chan string
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
	}

	return server
}

func (this *Server) handle(connect net.Conn) {
	fmt.Println("连接建立成功")

	user := NewUser(connect)

	user.Online()

	isLiveChannle := make(chan bool)

	go func() {
		buf := make([]byte, 4096)

		n, err := connect.Read(buf)

		if n == 0 {
			user.Offline()
			return
		}

		if err != nil && err != io.EOF {
			fmt.Println("connect Read err:", err)
			return
		}
		//去掉\n
		msg := buf[:n-1]

		user.DoMessage(string(msg))

		isLiveChannle <- true
	}()
	for {
		select {
		case <-isLiveChannle:
		case <-time.After(time.Second * 10):
			user.SendMsg("你被踢了")
			close(user.Channle)
			connect.Close()
			return
		}
	}
}

func (this *Server) BroadCast(user *User, msg string) {
	msg = "[" + user.Addr + "]" + user.Name + msg
	this.Channle <- msg
}

func (this *Server) ListenMessage() {
	for {
		msg := <-this.Channle

		this.Lock.Lock()
		for _, user := range this.Map {
			user.Channle <- msg
		}
		this.Lock.Unlock()
	}
}

func (this *Server) Start() {
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))

	if err != nil {
		fmt.Println("listen err:", err)
		return
	}

	defer listen.Close()

	go this.ListenMessage()

	for {
		connect, err := listen.Accept()

		if err != nil {
			fmt.Println("accept err:", err)
			continue
		}

		go this.handle(connect)
	}
}

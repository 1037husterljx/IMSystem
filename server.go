package main

import (
	"fmt"
	"net"
	"sync"
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
	//map放入用户
	this.Lock.Lock()
	this.Map[user.Name] = user
	this.Lock.Unlock()
	//广播上线
	this.BroadCast(user, "已上线")

	select {}
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

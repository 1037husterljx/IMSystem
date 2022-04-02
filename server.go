package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

func (this *Server) NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
	}

	return server
}

func (this *Server) handle(connect net.Conn) {
	fmt.Println("连接建立成功")
}

func (this *Server) Start() {
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("listen err:", err)
		return 
	}

	defer listen.Close()

	for {
		connect, err := listen.Accept()
		if err != nil {
			fmt.Println("accept err:", err)
			continue
		}
		go this.handle(connect)
	}
}

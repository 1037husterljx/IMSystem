package main

import "net"

type User struct {
	Name    string
	Addr    string
	connect net.Conn
	Channle chan string
}

func NewUser(connect net.Conn) *User {
	user := &User{
		Name:    connect.RemoteAddr().String(),
		Addr:    connect.RemoteAddr().String(),
		connect: connect,
		Channle: make(chan string),
	}

	//监听Channel
	go user.ListenChannle()

	return user
}

func (this *User) ListenChannle() {
	for {
		msg := <-this.Channle
		this.connect.Write([]byte(msg + "\n"))
	}
}

package main

import "net"

type User struct {
	Name    string
	Addr    string
	connect net.Conn
	Channle chan string
	server  *Server
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

func (this *User) Online() {
	//map放入用户
	this.server.Lock.Lock()
	this.server.Map[this.Name] = this
	this.server.Lock.Unlock()
	//广播上线
	this.server.BroadCast(this, "已上线")
}

func (this *User) Offline() {
	//map放入用户
	this.server.Lock.Lock()
	delete(this.server.Map,this.Name)
	this.server.Lock.Unlock()
	//广播上线
	this.server.BroadCast(this, "已下线")
}

func (this *User) DoMessage(msg string) {
	this.server.BroadCast(this,msg)
}

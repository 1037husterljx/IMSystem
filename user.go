package main

import (
	"net"
	"strings"
)

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
		this.SendMsg(msg)
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
	delete(this.server.Map, this.Name)
	this.server.Lock.Unlock()
	//广播上线
	this.server.BroadCast(this, "已下线")
}

func (this *User) DoMessage(msg string) {
	if msg == "who" {

		this.server.Lock.Lock()
		for _, user := range this.server.Map {
			this.SendMsg("[" + user.Name + "]" + "在线")
		}
		this.server.Lock.Unlock()

	} else if len(msg) > 7 && msg[:7] == "rename|" {
		newName := msg[7:]

		_, ok := this.server.Map[newName]

		if ok {
			this.SendMsg("用户名已被使用")
		} else {
			this.server.Lock.Lock()
			delete(this.server.Map, this.Name)
			this.server.Map[newName] = this
			this.server.Lock.Unlock()

			this.Name = newName
			this.SendMsg("已更新用户名" + newName)
		}

	} else if len(msg) > 4 && msg[:3] == "to|" {
		remoteName := strings.Split(msg, "|")[1]

		if remoteName == "" {
			this.DoMessage("请使用 to|name|消息 的格式")
		}

		remoteUser, ok := this.server.Map[remoteName]

		if !ok {
			this.SendMsg("用户不存在")
			return
		}

		content := strings.Split(msg, "|")[2]

		if content == "" {
			this.DoMessage("请使用 to|name|消息 的格式")
		}

		remoteUser.DoMessage(this.Name + "对你说" + content)
	} else {
		this.server.BroadCast(this, msg)
	}
}

func (this *User) SendMsg(msg string) {
	this.connect.Write([]byte(msg + "\n"))
}

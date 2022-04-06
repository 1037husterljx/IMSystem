<a name="NHxgv"></a>
# 系统架构
![image.png](https://cdn.nlark.com/yuque/0/2022/png/22805804/1649254614329-dbdaf342-f1bf-4074-a512-df32890cea00.png#clientId=uabfef3a1-a0ea-4&crop=0&crop=0&crop=1&crop=1&from=paste&height=395&id=ud0a9594e&margin=%5Bobject%20Object%5D&name=image.png&originHeight=978&originWidth=1683&originalType=binary&ratio=1&rotation=0&showTitle=false&size=706759&status=done&style=none&taskId=u613fae12-d963-49cc-8ded-14ad13519c1&title=&width=679.9999606970607)<br />Message用于广播，利用golang的channel进行收发信息，map用于保存<username，user>。
<a name="u1iIk"></a>
# 基础Server
tag:v1.0<br />利用net包进行socket操作
```go
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
```
<a name="vKP3x"></a>
# 用户上线及用户消息
tag：v2.0，3.0，4.0<br />信息放入message的channel中并广播给所有user的channel，对用户业务层进行封装，实现上线，下线，发送消息<br />**注意点**：user使用map保存,map需上锁
```go
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
```
<a name="THGh4"></a>
# 用户功能的升级
tag:v5.0，6.0<br />利用who指令对map中的用户进行查询，利用rename指令更改名字，利用to指令私聊，增加超时下线功能<br />**注意点**：以上功能不进行广播，直接在connect进行write
```go
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
```
```go
		user.DoMessage(string(msg))
        
		isLiveChannle <- true
	}()

	for {
        //利用了select关键字
		select {
		case <-isLiveChannle:
		case <-time.After(time.Second * 10):
			user.SendMsg("你被踢了")
			close(user.Channle)
			connect.Close()
			return
		}
	}
```
<a name="CDerK"></a>
# 客户端
tag:v7.0<br />利用net.Dial建立连接并封装上面的功能
```go
func main() {
    //go命令行的解析，在init函数可以获取参数用来获取server的位置
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

```
```go
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
```
```go
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
```
**注意点：**需要把connect数据打印在客户端
```go
go client.printRespose()

client.Run()
```
```go
func (this *Client) printRespose() {
	io.Copy(os.Stdout, this.connect)
}
```
<a name="BGElP"></a>
# 总结
此项目为学习golang而完成的一个简单即时通讯系统，主要利用了go的channel机制来完成一个即时通讯的功能，主要学习到了channel和协程的用法，并掌握了go的简单网络编程，并进行了多次版本迭代。在特性方面，有以下几点：

- 对协程和channel的使用
- defer关键字来关闭连接
- select关键字实现下线功能
- map及lock的使用
- 协程使用io.copy来阻塞等待网络响应


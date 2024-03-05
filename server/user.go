package server

import (
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

type User struct {
	Name    string
	Addr    string
	Channel chan string
	Conn    net.Conn
	Sever   *Server
	Live    chan bool
}

func NewUser(name string, conn net.Conn, server *Server) *User {
	user := &User{
		Name:    name,
		Addr:    conn.RemoteAddr().String(),
		Channel: make(chan string),
		Conn:    conn,
		Sever:   server,
		Live:    make(chan bool),
	}
	go user.ListenMessage()
	return user
}

func (u *User) ListenMessage() {
	for {
		msg := <-u.Channel
		u.Conn.Write([]byte(msg + "\n"))
	}
}

func (u *User) Online() {
	server := u.Sever
	server.MapLock.Lock()
	server.OnlineMap[u.Name] = u
	server.MapLock.Unlock()

	server.Broadcast(u, "欢迎加入聊天室")

}

func (u *User) Offline() {

	server := u.Sever

	server.MapLock.Lock()
	delete(server.OnlineMap, u.Name)
	server.MapLock.Unlock()

	u.Sever.Broadcast(u, "下线了")

	close(u.Channel)
	close(u.Live)

	u.Conn.Close()
}

func (u *User) ParseMsg(msg string) {
	server := u.Sever
	u.Live <- true
	if msg == "who" {
		server.MapLock.RLock()
		for _, v := range server.OnlineMap {
			onlineMsg := fmt.Sprintf("%s:%s 在线", v.Name, v.Addr)
			u.Channel <- onlineMsg
		}
		server.MapLock.RUnlock()
	} else if len(msg) > 7 && msg[0:7] == "rename " {
		server.MapLock.RLock()
		_, ok := server.OnlineMap[msg[7:]]
		if ok {
			u.Channel <- "改名失败，名字已存在"
		} else {
			delete(server.OnlineMap, u.Name)
			u.Name = msg[7:]
			server.OnlineMap[u.Name] = u
			u.Channel <- "改名成功"
		}
		server.MapLock.RUnlock()
	} else if len(msg) > 3 && msg[0:3] == "to " {
		userName := strings.Split(msg, " ")[1]
		if user, ok := server.OnlineMap[userName]; ok {
			user.Send(u.Name + ":" + msg[3+len(userName):])
		} else {
			u.Send("用户不存在")
		}
	} else {
		u.Sever.Broadcast(u, msg)
	}

}

func (u *User) HeartBeat() {
	for {
		select {
		case <-u.Live:

		case <-time.After(300 * time.Second):
			u.Send("你被下线了")
			return
		}
	}
}

func (u *User) Send(msg string) {
	u.Channel <- msg
}

func (u *User) ListenUserMessage() {
	buf := make([]byte, 1024)
	conn := u.Conn
	var msg string
	for {
		n, err := conn.Read(buf)

		if err != nil && err != io.EOF {
			fmt.Println("conn.read err:", err)
			return
		}

		if n == 0 {
			u.Offline()
			return
		} else {
			msg = string(buf[:n])
			if strings.Contains(msg, "\n") {
				msg = msg[:len(msg)-1]
			}
		}
		fmt.Println(msg)
		u.ParseMsg(msg)
	}
}

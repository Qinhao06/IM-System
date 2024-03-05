package server

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	OnlineMap map[string]*User
	MapLock   sync.RWMutex

	BroadcastMessage chan string
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:               ip,
		Port:             port,
		OnlineMap:        make(map[string]*User),
		BroadcastMessage: make(chan string),
	}
	return server
}

func (server *Server) Broadcast(user *User, msg string) {
	sendMessage := "[" + user.Addr + "]" + user.Name + ": " + msg
	server.BroadcastMessage <- sendMessage
}

func (server *Server) ListenMessage() {
	for {
		msg := <-server.BroadcastMessage
		server.MapLock.Lock()
		for _, user := range server.OnlineMap {
			user.Channel <- msg
		}
		server.MapLock.Unlock()
	}
}

func (server *Server) Handler(conn net.Conn) {
	user := NewUser(conn.RemoteAddr().String(), conn, server)
	user.Online()
	go user.ListenUserMessage()
	go user.HeartBeat()
}

func (server *Server) Start() {

	//socket listen
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.Ip, server.Port))
	if err != nil {
		fmt.Println("net.listen err:", err)
		return
	}

	// socket.close
	defer listen.Close()

	go server.ListenMessage()
	//accept
	for {
		accept, err := listen.Accept()
		if err != nil {
			fmt.Print("listen accept err:", err)
			continue
		}

		server.Handler(accept)

	}
}

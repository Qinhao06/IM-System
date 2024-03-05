package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	Conn       net.Conn
	Flag       int
}

func NewClient(ip string, port int, name string) *Client {
	client := &Client{
		ServerIp:   ip,
		ServerPort: port,
		Flag:       999,
	}

	dial, err := net.Dial("tcp", ip+":"+strconv.Itoa(port))
	if err != nil {
		fmt.Println("连接服务器失败")
		return nil
	}
	client.Conn = dial
	return client
}

func (c *Client) Send(msg string) {
	c.Conn.Write([]byte(msg))
}

func (c *Client) Listen() {
	for {
		data := make([]byte, 1024)
		n, err := c.Conn.Read(data)
		if err != nil {
			fmt.Println("连接服务器失败")
			panic(err)
		}
		fmt.Print("服务器发来消息 > " + string(data[:n]))
	}
}

func (c *Client) Start() {
	go c.Listen()

	for c.Flag != 0 {
		c.Menu()
		switch c.Flag {
		case 1:
			c.changeName()
		case 2:
			c.privateChat()
		case 3:
			c.publicChat()
		}
	}

}

func (c *Client) readStdin() string {
	reader := bufio.NewReader(os.Stdin)
	var msg string
	msgData, _, _ := reader.ReadLine()
	msg = string(msgData)
	if msg == "" {
		return c.readStdin()
	}
	return msg
}

func (c *Client) changeName() {
	name := ""
	for name == "" {
		fmt.Println("请输入用户名")
		name = c.readStdin()
	}
	c.Send("rename " + name)
}

func (c *Client) privateChat() {
	c.Send("who")
	var remoteName string
	fmt.Println("请选择私聊对象，并输入用户名")
	fmt.Scanln(&remoteName)
	fmt.Println("请输入消息, 输入 exit 退出")
	for {
		msg := c.readStdin()
		if msg == "exit" {
			break
		}
		c.Send("to " + remoteName + " " + msg)
	}

}

func (c *Client) publicChat() {
	fmt.Println("请输入消息, 输入 exit 退出")
	for {
		msg := c.readStdin()
		if msg == "exit" {
			break
		}
		c.Send(msg)
	}
}

func (c *Client) Menu() {
	fmt.Println("输入数字选择模式")
	fmt.Println("0 退出")
	fmt.Println("1 更新用户名")
	fmt.Println("2 私聊模式")
	fmt.Println("3 群聊模式")
	fmt.Scanln(&c.Flag)
	if !(c.Flag >= 0 && c.Flag <= 3) {
		fmt.Println("输入错误,请重新输入")
		c.Menu()
	}
}

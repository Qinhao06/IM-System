package main

import (
	client2 "IM-System/client"
	"fmt"
	"strings"
)

var serverIp string = "127.0.0.1"
var serverPort int = 8888
var name string = "client"

func parse() {
	fmt.Println("请输入服务器ip")
	fmt.Scanf("%s", &serverIp)
	if len(strings.Split(serverIp, ".")) != 4 {
		fmt.Println("ip 输入不正确，使用默认 ip 127.0.0.1")
	}

	fmt.Println("请输入端口号")
	fmt.Scanf("%d", &serverPort)
	fmt.Println("请输入Name:")
	fmt.Scanf("%s", &name)
}

func main() {
	parse()
	client := client2.NewClient(serverIp, serverPort, name)
	if client == nil {
		fmt.Print("创建客户端失败")
		return
	}
	fmt.Println("链接服务器成功")
	client.Start()
}

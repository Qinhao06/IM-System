package main

import (
	"IM-System/server"
	"fmt"
)

func main() {
	server := server.NewServer("127.0.0.1", 8888)
	server.Start()
	fmt.Print("sever closed")
}

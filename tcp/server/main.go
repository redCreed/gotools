/*
   文件名：server.go
   server服务端的示例代码（未处理粘包问题）
   服务端接收到数据后立即打印
   此时将会不间断的出现TCP粘包问题
*/
package main

import (
	"fmt"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":9000")
	if err != nil {
		return
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go handleConnection(conn)
	}
}
func handleConnection(conn net.Conn) {
	defer conn.Close()
	for {
		buf := make([]byte, 1024)
		_, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Read Error:", err)
			return
		}
		fmt.Println(string(buf) + "\r\n")
	}
}

package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func handle(conn net.Conn, logger *log.Logger) { //处理连接方法
	defer conn.Close() //关闭连接
	for {
		buf := make([]byte, 100)
		n, err := conn.Read(buf) //读取客户端数据
		if err != nil {
			fmt.Println(err)
			return
		}

		logger.Println(string(buf[0:n]))
	}
}

func main() {
	fmt.Println("start coordinator...")
	listen, err := net.Listen("tcp", "0.0.0.0:11300") //创建监听
	if err != nil {
		fmt.Println("listen failed! msg :", err)
		return
	}

	// init log
	file, err := os.OpenFile("coordinator.txt", os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Fatalf("create log failed: %v\n", err)
	}

	logger := log.New(file, "", log.Lshortfile|log.LstdFlags)

	for {
		conn, errs := listen.Accept() //接受客户端连接
		if errs != nil {
			fmt.Println("accept failed")
			continue
		}
		go handle(conn, logger) //处理连接
	}
}

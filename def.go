package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const (
	HEADER = "250"//数据包开始标志
	PAYLOAD_SIZE = 1//数据长度占用字节,这里其实只保存20这个值,所以1个字节就够了
)

/**
 * 服务器A返回的数据
 * 或者服务器B经过服务器A转发回来的数据
 */
func readSomething(conn *net.Conn,more bool) []byte {
	buffer := make([]byte, 12)
	for{
		_, err := (*conn).Read(buffer)
		if isErrAPrint(err){
			//连接出现问题或服务器关闭连接,退出.
			panic(fmt.Sprintf("Read error: %s", err))
		}
		//这里处理服务器返回的服务数据
		fmt.Println(string(buffer))

		if !more {return nil}
	}
}

func writeSomething(conn *net.Conn, data []byte) bool {
	_, err := (*conn).Write(data)
	if isErrAPrint(err) {
		//连接出现问题,退出.
		//panic(fmt.Sprintf("Write error: %s", err))
		return false
	}
	return true
}

func isErrAPrint(err error) bool {
	if err != nil{
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		return true
	}
	return false
}

func watch() {
	sigs := make(chan os.Signal)
	signal.Notify(sigs,syscall.SIGINT)
	select {
	case sig := <- sigs:
		if sig == syscall.SIGINT{
			os.Exit(0)
		}
	}
}

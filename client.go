package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"gopkg.in/snksoft/crc.v1"
	"net"
	"os"
	"strconv"
	"time"
)

type Client struct {
	Sip string
	Sport uint64
	ID string
}

func (client *Client)Run()  {
	conn, err := net.DialTimeout("tcp",client.Sip+":"+strconv.FormatUint(client.Sport, 10),3*time.Second)
	if isErrAPrint(err) { os.Exit(1) }
	defer conn.Close()

	for i:=0; i<20; i++{
		if client.login(&conn) {
			break
		}

		if i == 20 {
			fmt.Println("20次登陆失败,不再登陆.")
			os.Exit(1)
		}
		time.Sleep(time.Second)
	}

	client.todoSomething(&conn)
}

func (client *Client)login(conn *net.Conn) bool {
	payloadLen := len(client.ID)
	var buffer bytes.Buffer

	var ck = make([]byte, 2)
	hash := crc.NewHash(crc.X25)
	x25Crc := hash.CalculateCRC([]byte(client.ID))
	binary.LittleEndian.PutUint16(ck, uint16(x25Crc))

	payloadLenbytes := make([]byte,PAYLOAD_SIZE)
	payloadLenbytes = []byte{byte(payloadLen)}

	buffer.Write([]byte(HEADER))
	buffer.Write(payloadLenbytes)
	buffer.WriteString(client.ID)
	buffer.Write(ck)

	data := buffer.Bytes()
	isLogin := true
	defer func() {
		if err:=recover();err!=nil{
			isLogin = false
		}
	}()
	writeSomething(conn,data)
	if isLogin {readSomething(conn)}//如果这个函数没有panic,就说明登陆成功了,登陆失败服务器会关闭连接,报panic.
	return isLogin
}

func (client *Client)todoSomething(conn *net.Conn)  {
	go readSomething(conn)
	for{
		writeSomething(conn,[]byte("abc"))
		time.Sleep(time.Second)
	}
}
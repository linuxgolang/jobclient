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

		if i == 19 {
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
	writeSomething(conn,data)
	isLogin := true
	defer func() {
		if err:=recover();err!=nil{
			isLogin = false
		}
	}()
	readSomething(conn,false)
	fmt.Println("ccccccc")
	return isLogin
}

func (client *Client)todoSomething(conn *net.Conn)  {
	go func() {
		for{
			if !writeSomething(conn,[]byte("abc")) {return}
			time.Sleep(time.Second)
		}
	}()
	readSomething(conn,true)
}
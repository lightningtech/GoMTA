package mta

import (
	"log"
	"net"
	//"encoding/hex"
	"io"
)

type Mta struct{
	Host string
}

var currentMessageBody MessageBody
var currentHost string

func (m Mta) Send(messageBody MessageBody) {
	currentHost = m.Host
	send(messageBody)
}

func send(messageBody MessageBody) {
	currentMessageBody = messageBody
	sendMessage(currentHost, messageBody)
}

func retrySend() {
	log.Println("retrying sending of message")
	send(currentMessageBody)
}

func sendMessage(host string, messageBody MessageBody) {
	conn, err := net.Dial("tcp", host + ":smtp")
	if err != nil {
		log.Fatal("err ", err)
	}
	onResponse(conn)
	for i, command := range messageBody.Data {
		sendCommand(conn, command, i <= messageBody.DataCommandIndex || i >= messageBody.EndDataCommandIndex)
	}
}

func sendCommand(conn net.Conn, command []byte, hasResponse bool) {
	conn.Write(command)
	if hasResponse {
		onResponse(conn)
	}
}

func onResponse(conn net.Conn) {
	data := make([]byte, 1024)
	_, err := conn.Read(data)
	if err != nil {
		if err != io.EOF {
			panic(err)
		} else {
			log.Printf("err: %v", err)
		}
		retrySend()
	}
}
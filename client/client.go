package main

import (
	"net"
	"os"
	"fmt"
	"strings"
	"bufio"
	"encoding/json"
	"io"
)

const (
	HOST = "192.168.6.138"
	PORT = "8080"
	TYPE = "tcp"
)

type Message struct {
	Username string `json:"username"`
	Message string `josn:"Message"`
}

func main() {
	tcpServer , err := net.ResolveTCPAddr(TYPE, HOST+":"+PORT)
	if err != nil {
		println("ResolveTCPAddr failed", err.Error())
		os.Exit(1)
	}
	
	conn, err := net.DialTCP(TYPE, nil, tcpServer)
	stop := make(chan bool)
	if err != nil {
		println("Dial failed: ", err.Error())
		os.Exit(1)
	}
	defer conn.Close()	
	username := sendUsername(conn)

	fmt.Print("You: ")
	go handleMessages(conn, stop)
	go handleUserInput(username, conn, stop)

	<-stop
	fmt.Println("Closing connection")
}

func handleMessages(conn net.Conn, stop chan<- bool) {
	for {
		var m Message
		decoder := json.NewDecoder(conn)
		if err:= decoder.Decode(&m); err != nil {
			if err == io.EOF {
				fmt.Println("Server unavailable")
			} else {
				fmt.Println("Unable to decode the message", err)
			}
			os.Exit(1)
		}
		fmt.Print("\r           \r")
		fmt.Print(m.Username + ": " + m.Message)
		fmt.Print("You: ")
	}
}

func handleUserInput(username string, conn net.Conn, stop chan<-bool) {
	reader := bufio.NewReader(os.Stdin)
	response := Message{Username: username}
	for {
		msg, _ := reader.ReadString('\n')
		if strings.TrimSpace(msg) == "" {
			continue
		}
		if msg == "exit" {
			stop <- true
		}

		fmt.Print("You: ")
		response.Message = msg
		json, err := json.Marshal(response)	
		if err != nil {
			fmt.Println("Unable to marshal the response")
			continue
		}
		conn.Write(json)
	}
}

func sendUsername(conn net.Conn) string {
		received := make([]byte, 1024)
		n, err := conn.Read(received)
		if err != nil {
			println("Read data failed:", err.Error())
			os.Exit(1)
		}
		fmt.Println(string(received[:n]))
		var username string
		fmt.Scan(&username)
		conn.Write([]byte(username))
		return username
}

var Esc = "\x1b"

func moveUp(n int) string {
	return escape("[%dA", n)
}

func escape(format string, args ...interface{}) string {
	return fmt.Sprintf("%s%s", Esc, fmt.Sprintf(format, args...))
}

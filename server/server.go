package main

import (
	"net"
	"log"
	"os"
	"io"
	"sync"
	"strings"
)

const (
	HOST = "192.168.6.138"
	PORT = "8080"
	TYPE = "tcp"
)

type Username = string

type message struct {
	user string
	message string
}


var clients = struct {
	sync.Mutex
	conns map[string]net.Conn
}{conns: make(map[Username]net.Conn)}

func main() {
	log.Println("Starting server")
	listner, err := net.Listen(TYPE, HOST+":"+PORT)
	if err != nil {
		log.Fatal("Error creating a listner: ", err)
		os.Exit(1)
	}
	log.Printf("Server started at: %s:%s", HOST, PORT)
	defer listner.Close()

	for {
		conn, err := listner.Accept()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		log.Printf("Connection established: %s", conn.RemoteAddr().String())
		go handleRequest(conn)
	}
}

func getConnectedUsers() []string{
	clients.Lock()
	users := make([]string, 0, len(clients.conns))
	for key := range clients.conns {
		users = append(users, key)
	}
	clients.Unlock()
	return users
}

func handleRequest(conn net.Conn) {
		conn.Write([]byte("Enter your username: "))
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			log.Printf("Error occured while getting username: %s", err.Error())
			conn.Close()
		}
		username := string(buffer[:n])
		clients.Lock()
		clients.conns[username] = conn
		clients.Unlock()
		log.Println(getConnectedUsers())
		defer func() {
			clients.Lock()
			delete(clients.conns, username)
			clients.Unlock()
			conn.Close()
			log.Printf("Client disconnected: %s", username)
		}()
			
		for {
			n, err := conn.Read(buffer)
			if err != nil {
				if err == io.EOF {
					log.Printf("Closing connection: %s", conn.RemoteAddr().String())
				} else {
					log.Printf("Error reading the message %s", err)
				}
				break
			}
			message := strings.TrimSpace(string(buffer[:n]))
			if message == "" {
				log.Println("Received an empty message from: ", username)
				continue
			}
			log.Printf("Message received from %s", username)
			for _, v := range clients.conns {
				if v == conn{
					continue
				}
				v.Write(buffer[:n])
			}
		}
}


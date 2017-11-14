//use tcp realize broadcast message
package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

var buffServer = make([]byte, 1024)
var buffClient = make([]byte, 1024)
var clients = make(map[string]net.Conn)
var messages = make(chan string, 10)

func TcpServer(port string) {
	port = ":" + port
	tcpAddr, err := net.ResolveTCPAddr("tcp4", port)
	if err != nil {
		log.Printf("error:%s\n", err)
		return
	}
	listenTcp, err := net.ListenTCP("tcp", tcpAddr)
	defer listenTcp.Close()
	if err != nil {
		log.Printf("error:%s\n", err)
		return
	}
	log.Printf("开启了服务端:%s\n", "正在监听客户端的请求")
	go BoradCast(clients, messages)
	for {
		conn, err := listenTcp.Accept()
		if err != nil {
			log.Printf("error:%s\n", err)
		}
		clients[conn.RemoteAddr().String()] = conn
		go HandleServer(conn, messages)
	}
}

func HandleServer(conn net.Conn, messages chan string) {
	for {
		_, err := conn.Read(buffServer)
		if err != nil {
			log.Printf("error:%s\n", err)
			conn.Close()
			return
		}
		messages <- string(buffServer)
	}
}

func BoradCast(clients map[string]net.Conn, messages chan string) {
	for {
		msg := <-messages
		for index, client := range clients {
			_, err := client.Write([]byte(msg))
			if err != nil {
				log.Printf("error:%s\n", err)
				delete(clients, index)
			}
		}
	}
}

func TcpClient(serverAddr string) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", serverAddr)
	if err != nil {
		log.Printf("error:%s\n", err)
		return
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	defer conn.Close()
	if err != nil {
		log.Printf("error:%s\n", err)
		return
	}
	log.Printf("客户端:%s开始连接\n", conn.LocalAddr().String())
	go ClientSend(conn)
	for {
		_, err := conn.Read(buffClient)
		if err != nil {
			log.Printf("error:%s\n", err)
			return
		}
		log.Printf("收到消息:%s\n", string(buffClient))
	}
}

func ClientSend(conn net.Conn) {
	var input string
	username := conn.LocalAddr().String()
	for {
		fmt.Scanln(&input)
		if input == "/quit" {
			log.Printf("Info:%s\n", "ByeBye")
			conn.Close()
			os.Exit(0)
		}
		_, err := conn.Write([]byte(username + ":" + input))
		if err != nil {
			log.Printf("error:%s\n", err)
			return
		}
	}
}

//开启TcpServer: go run main.go server 端口号(8080)
//开启TcpClient: go run main.go client 本地ip地址:端口号(8080)
func main() {
	if len(os.Args) != 3 {
		log.Println("command error")
	}
	if os.Args[1] == "server" {
		TcpServer(os.Args[2])
	}
	if os.Args[1] == "client" {
		TcpClient(os.Args[2])
	}
}

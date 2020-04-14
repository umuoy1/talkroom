package main

import (
	"fmt"
	"log"
	"net"
	"sync"
)

type Client struct {
	addr string
	C    chan string
}

func main() {
	fmt.Println("Server is runing on 127.0.0.1 at port 3000")
	listenner, err := net.Listen("tcp", "127.0.0.1:3000")
	defer listenner.Close()
	if err != nil {
		log.Println(err)
		return
	}
	go router()
	for {
		conn, err := listenner.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		//处理客户端连接
		go handleClient(conn)
	}
}

func msgToClient(conn net.Conn, cli Client) {
	for ms := range cli.C {
		conn.Write([]byte(ms))
	}
}

var onlineMap = make(map[string]Client)
var rLock sync.RWMutex
func handleClient(conn net.Conn) {
	fmt.Println(conn.RemoteAddr().String(), "接入")
	msg <- conn.RemoteAddr().String()+" 接入"
	cli := Client{
		addr: conn.RemoteAddr().String(),
		C:    make(chan string),
	}
	rLock.Lock()
	onlineMap[conn.RemoteAddr().String()] = cli
	rLock.Unlock()
	//用于向当前客户端发送消息
	go msgToClient(conn, cli)
	//处理客户端发送的消息
	go execData(conn)
}

var msg = make(chan string)

//消息转发
func router() {
	for {
		ms := <-msg
		//给每个在线用户发消息
		rLock.RLock()
		for _, cli := range onlineMap {
			cli.C <- ms
		}
		rLock.RUnlock()
	}
}

//处理客户端发送的数据
func execData(conn net.Conn) {
	for {
		var data = make([]byte, 1024*4)
		n, err := conn.Read(data)
		if err != nil || string(data[:n]) == "bye"{
			fmt.Println(conn.RemoteAddr().String(), "断开")
			msg<-conn.RemoteAddr().String()+" 断开"
			rLock.Lock()
			delete(onlineMap, conn.RemoteAddr().String())
			rLock.Unlock()
			return
		}
		fmt.Println(conn.RemoteAddr().String(), "发送：", string(data[:n]))
		//消息放到转发器中
		msg <- conn.RemoteAddr().String() + ": " + string(data[:n])
	}
}
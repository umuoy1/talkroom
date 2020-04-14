package main

import (
	"fmt"
	"net"
	"log"
)
func main() {
	cl,err := net.Dial("tcp", "127.0.0.1:3000")
	if err != nil{
		log.Println(err)
	}
	go func() {
		for {
			var data = make([]byte, 1024 * 4)
			n,err := cl.Read(data)
			if err!=nil{
				continue
			}
			fmt.Println(string(data[:n]))
		}
	}()
	for {
		var data = make([]byte, 1024 * 4)
		fmt.Scan(&data)
		if string(data)=="exit"{
			cl.Close()
			return
		}
		_,err := cl.Write(data)
		if err!=nil{
			return
		}
	}
	cl.Close()
}
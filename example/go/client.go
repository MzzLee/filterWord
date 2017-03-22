package main

import (
	"net"
	"fmt"
	"os"
	"io"
	"encoding/json"
	"strconv"
	"time"
)

var (
	HeaderPrefix = "[HDR]"
	HeaderSuffix = "[/HDR]"
)

func main() {
	start := time.Now().Unix()
	conn := Connect("127.0.0.1", 8821)
	for i:=0;i<10;i++ {
		Send(conn, "nimeiafuck消息队列等待你888流浪者" + strconv.Itoa(i))
		//fmt.Println("send success!")
		//time.Sleep(1 * time.Second)
	}
	end := time.Now().Unix()
	fmt.Println(end-start)
	time.Sleep(1 * time.Second)


}

func Connect(address string, port int) net.Conn{
	conn, err := net.Dial("tcp", address + ":" + strconv.Itoa(port))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
	return conn
}

func Pack(body string, keepAlive int) string{
	header := make(map[string]int)
	header["content-length"] = len(body)
	header["keep-alive"] = keepAlive
	jsonHeader, err :=json.Marshal(header)
	if err != nil{
		fmt.Println("Json encode Error : ", err.Error())
		return HeaderPrefix + "0" + HeaderSuffix
	}
	return HeaderPrefix + string(jsonHeader) + HeaderSuffix + body
}

func Send (conn net.Conn,  content string) {

	words :=Pack(content, 1)
	io.WriteString(conn, words)
	data  := make([]byte, 4096)
	tmpBuffer  := make([]byte, 0)
	i := 0
	for {
		n, err := conn.Read(data)
		if err != nil {
			break
		}
		tmpBuffer = append(tmpBuffer, data[:n]...)
		for _, v := range tmpBuffer{
			if v != 13{
				i++
			}else{
				fmt.Println(string(tmpBuffer[0:i]))
				return
			}
		}
	}

	//conn.Close()
}
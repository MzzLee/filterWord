package main

import (
	"net"
	"fmt"
	"os"
	"io"
	"encoding/json"
	"time"
	"strconv"
)

var (
	HeaderPrefix = "[header]"
	HeaderSuffix = "[/header]"
)

func main() {

	for i:=0;i<50000;i++{
		conn := Connect("127.0.0.1", 8821)
		go Send(conn, "fuck day ! ")
	}
	time.Sleep(time.Second * 1)

}

func Connect(address string, port int) net.Conn{
	conn, err := net.Dial("tcp", address + ":" + strconv.Itoa(port))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
	return conn
}

func Pack(body string, isAlive int) string{
	jsonData, err := json.Marshal(body)

	header := make(map[string]int)
	header["content-length"] = len(string(jsonData))
	header["is-alive"] = isAlive
	jsonHeader, err :=json.Marshal(header)
	if err != nil{
		fmt.Println("Json encode Error : ", err.Error())
		return HeaderPrefix + "0" + HeaderSuffix
	}
	return HeaderPrefix + string(jsonHeader) + HeaderSuffix + string(jsonData)
}

func Send (conn net.Conn,  content string) {

	words :=Pack(content, 0)
	io.WriteString(conn, words)
	//start := time.Now().Unix()
	var data  = make([]byte, 4096)
	for {
		count, err := conn.Read(data)
		if err != nil {
			break
		}
		fmt.Println(string(data[:count]))
	}
	conn.Close()
	//end := time.Now().Unix()
	//fmt.Println(end-start)
}
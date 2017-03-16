package main

import (
	"net"
	"fmt"
	"os"
	"io"
	"encoding/json"
	"time"
)

var (
	HeaderPrefix = "[header]"
	HeaderSuffix = "[/header]"
)

func main() {

	conn, err := net.Dial("tcp", "127.0.0.1:9901")

	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	fmt.Println("content success")
	for i:=0;i<50000;i++{
		go Send(conn)
	}
	time.Sleep(time.Second * 1)
	conn.Close()
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

func Send (conn net.Conn) {

	words :=Pack("fuck你妹的插暴王嘉龙对妹子抓胸triangleisteday", 0)
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
	//end := time.Now().Unix()
	//fmt.Println(end-start)
}
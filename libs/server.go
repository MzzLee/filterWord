package libs

import (
	"net"
	"log"
	"os"
	"strings"
	"strconv"
	"encoding/json"
	"errors"
	"fmt"
)
const (
	StatusOK = 200
	StatusRequestError = 400
)

const (
	HeaderPrefix = "[Header]"
	HeaderPrefixLen = 8
	HeaderSuffix = "[/Header]"
	HeaderSuffixLen = 9
	HeaderMLen   = 1024
)

type Response struct{
	Status uint16 `json:"status"`
	Body string `json:"body"`
}

type Server struct {
	Ac *Node
	Host string
	Port int
	Protocol string
}


var _serverInstance *Server

func Create(ac *Node, config *Conf) *Server{
	if _serverInstance == nil {
		_serverInstance = new(Server)
		_serverInstance.Ac = ac
		_serverInstance.Host = config.Bind
		_serverInstance.Port = config.Port
		_serverInstance.Protocol = config.Protocol
	}

	return _serverInstance
}

func (server *Server)Start(){
	listenFD, err := net.Listen(server.Protocol, server.Host + ":" + strconv.Itoa(server.Port))
	if err !=nil {
		log.Fatalf("Error %s \n", err.Error())
		os.Exit(1)
	}

	fmt.Println("Server start successed !")

	defer listenFD.Close()

	for {
		conn, err := listenFD.Accept()
		if err != nil {
			continue
		}
		go server.receive(conn)
	}
}

func (server *Server) receive (conn net.Conn){
	tmpBuf := make([]byte, 0)

	defer conn.Close()

	responseStatus := 0
	buffer := make([]byte, HeaderMLen)
	requestLength := 0
	runTimes :=0

	for{
		n, err := conn.Read(buffer)

		if err != nil {
			log.Fatalln("conn closed")
			return
		}
		tmpBuf = append(tmpBuf, buffer[:n]...)
		//取响应头
		if requestLength == 0 && strings.HasPrefix(string(tmpBuf), HeaderPrefix){
			if suffix := strings.Index(string(tmpBuf), HeaderSuffix); suffix > 0 {
				requestLength, _ = strconv.Atoi(string(tmpBuf[HeaderPrefixLen:suffix]))
				if requestLength == 0{
					responseStatus = StatusRequestError
					break
				}
				tmpBuf = tmpBuf[suffix + HeaderSuffixLen:]
			}

		}else{
			if runTimes == 0{
				responseStatus = StatusRequestError
				break
			}
		}

		runTimes++

		if requestLength !=0 && len(tmpBuf) >= requestLength{
			responseStatus = StatusOK
			break
		}
	}
	var responseBody string
	var err error
	if responseStatus == StatusOK{
		responseBody, err = server.Work(tmpBuf)
		if err != nil{
			responseStatus = StatusRequestError
		}
	}

	server.Response(conn, uint16(responseStatus), responseBody)
	return
}


func (server *Server) Response (conn net.Conn, status uint16, body string){
	var response = Response{status, body}
	responseResult, err := json.Marshal(response)
	if err !=nil{
		log.Fatal("Response Result Json Error")
	}
	conn.Write(responseResult)
}

func (server *Server) Work (request []byte) (string, error){
	defer func(){
		if err :=recover(); err !=nil {
			log.Fatalf("Work Error: %s", err)
		}
	}()
	var body string
	err := json.Unmarshal(request, &body)
	if err == nil{
		return server.Ac.AcFind(body), nil
	}
	return "", errors.New("Ac Found error")
}


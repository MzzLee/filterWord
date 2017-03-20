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
	"io"
)
const (
	StatusOK = 200
	StatusRequestError = 400
)

const (
	DefaultPort 	= 9901
	DefaultBind  	= "127.0.0.1"
	HeaderPrefix 	= "[header]"
	HeaderPrefixLen = 8
	HeaderSuffix	= "[/header]"
	HeaderSuffixLen = 9
	ReaderMLen   	= 2048
)

type Request struct {
	ContentLength int `json:"content-length"`
	IsAlive int `json:"is-alive"`
}

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

		if config.Port == 0 {
			config.Port = DefaultPort
		}

		if config.Bind == ""{
			config.Bind = _serverInstance.getRealIP()
		}

		_serverInstance.Ac = ac
		_serverInstance.Host = config.Bind
		_serverInstance.Port = config.Port
		_serverInstance.Protocol = config.Protocol
	}

	return _serverInstance
}

func (server *Server) getRealIP() string{
	ips, err :=  net.InterfaceAddrs()
	if err != nil{
		log.Fatalf("Error %s \n", err.Error())
		os.Exit(1)
	}

	for _, address := range ips {
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}
	return DefaultBind
}

func (server *Server) Start (){
	listenFD, err := net.Listen(server.Protocol, server.Host + ":" + strconv.Itoa(server.Port))
	if err !=nil {
		log.Fatalf("Error %s \n", err.Error())
		os.Exit(1)
	}

	fmt.Println("Server start successed! You can \"Ctrl + c\" quit!")

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
	tmpBuffer := make([]byte, 0)

	defer conn.Close()

	readBuffer := make([]byte, ReaderMLen)
	request    := Request{0,0}


	for {
		n, err := conn.Read(readBuffer)
		if err != nil || err == io.EOF {
			return
		}

		tmpBuffer = append(tmpBuffer, readBuffer[:n]...)

		var sendBuffer []byte
		for {
			if len(tmpBuffer) == 0 {
				break
			}

			if request.ContentLength == 0{
				tmpBuffer, request = server.getHeader(tmpBuffer)
				if request.ContentLength == 0 {
					server.Response(conn, uint16(StatusRequestError), "")
					return
				}
			}

			if request.ContentLength > 0 && len(tmpBuffer) >= request.ContentLength{

				responseStatus := StatusOK
				sendBuffer = tmpBuffer[0:request.ContentLength]

				tmpBuffer  = tmpBuffer[request.ContentLength:]
				responseBody, err := server.Work(sendBuffer)
				if err != nil {
					responseStatus = StatusRequestError
				}

				server.Response(conn, uint16(responseStatus), responseBody)
				if request.IsAlive > 0 {
					request.ContentLength = 0
				}else{
					return
				}

			}else{
				break
			}

		}
	}
	return
}

func (server *Server) getHeader (tmpBuffer []byte) ([]byte, Request){
	var request Request
	if strings.HasPrefix(string(tmpBuffer), HeaderPrefix) {
		suffix := strings.Index(string(tmpBuffer), HeaderSuffix)
		if suffix > 0 {
			json.Unmarshal(tmpBuffer[HeaderPrefixLen:suffix], &request)
			tmpBuffer = tmpBuffer[suffix+HeaderSuffixLen:]
		}
	}
	return tmpBuffer, request
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


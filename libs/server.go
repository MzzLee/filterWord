package libs

import (
	"net"
	"os"
	"strings"
	"strconv"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)
const (
	StatusOK = 200
	StatusRequestError = 400
)

const (
	DefaultPort 	= 9901
	DefaultBind  	= "127.0.0.1"
	HeaderPrefix 	= "[HDR]"
	HeaderPrefixLen = 5
	HeaderSuffix	= "[/HDR]"
	HeaderSuffixLen = 6
	ReaderMLen   	= 1024
	ResponseEOF     = byte(13)
)

type Package struct {
	Request *Request
	Content []byte
}

type Request struct {
	ContentLength int `json:"content-length"`
	KeepAlive int `json:"keep-alive"`
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

var (
	_serverInstance *Server
	_logger *GLogger
)

func CreateServer() *Server{
	if _serverInstance == nil {
		_serverInstance = new(Server)
		_serverInstance.Init()
	}

	return _serverInstance
}

func (server *Server) getRealIP () string{
	ips, err :=  net.InterfaceAddrs()
	if err != nil{
		_logger.Fatal(err.Error())
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

func (server *Server) Init () *Server{
	config := GetConfigInstance()
	server.Ac = AcBuild(config.Keyword)
	if config.Port == 0 {
		config.Port = DefaultPort
	}

	if config.Bind == ""{
		config.Bind = _serverInstance.getRealIP()
	}

	server.Host = config.Bind
	server.Port = config.Port
	server.Protocol = config.Protocol

	_logger, _ = InitLogger(config.LogFile, config.Env)
	return server
}

func (server *Server) Start (){
	listenFD, err := net.Listen(server.Protocol, server.Host + ":" + strconv.Itoa(server.Port))
	if err !=nil {
		_logger.Fatal(err.Error())
		fmt.Println(err.Error())
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

func (server *Server) receive (conn net.Conn) {

	tmpBuffer 	:= make([]byte, 0)
	readBuffer 	:= make([]byte, ReaderMLen)
	workChannel	:= make(chan []byte)
	request    	:= Request{0,0}
	headerStatus 	:= false
	pack   		:= Package{}
	isAlive		:= 0

	defer func(){
		time.Sleep(time.Millisecond)
		conn.Close()
	}()

	for {
		n, err := conn.Read(readBuffer)
		if err != nil {
			return
		}

		tmpBuffer = append(tmpBuffer, readBuffer[:n]...)
		for {
			if len(tmpBuffer) == 0 {
				break
			}

			go server.Work(conn, workChannel)

			//获取头文件
			if headerStatus == false {
				tmpBuffer, request, err = server.getHeader(tmpBuffer)
				pack.Request = &request
				if err != nil {
					_logger.Warning("Request Header Error : ", err.Error())
					if isAlive == 0{
						workBuffer,_ := json.Marshal(pack)
						workChannel <- workBuffer
						return
					}else{
						break
					}
				}else{
					isAlive = request.KeepAlive
					headerStatus = true
				}
			}

			//获取内容体
			if headerStatus {
				if request.ContentLength <= 0 {
					workBuffer,_ := json.Marshal(pack)
					workChannel <- workBuffer
					if request.KeepAlive == 0 {
						return
					}else{
						headerStatus = false
					}
				}else{
					if len(tmpBuffer) >= request.ContentLength {

						pack.Content = tmpBuffer[0:request.ContentLength]
						tmpBuffer    = tmpBuffer[request.ContentLength:]
						workBuffer,_ := json.Marshal(pack)
						workChannel <- workBuffer

						if request.KeepAlive == 0 {
							return
						}else{
							headerStatus = false


						}
					}else{
						break
					}
				}
			}else{
				break
			}

		}
	}

	return
}

func (server *Server) HeartBeating(conn net.Conn, readerChannel chan []byte, timeout int) {
	select {
	case fk := <-readerChannel:
		_logger.Trace("Request Heart Beating from ",conn.RemoteAddr().String(), string(fk))
		conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
		break
	case <-time.After(time.Second * 5):
		conn.Close()
	}
}

func (server *Server) GravelChannel(readBuffer []byte, readerChannel chan []byte) {
	readerChannel <- readBuffer
}

func (server *Server) getHeader (tmpBuffer []byte) ([]byte, Request, error){
	var request Request
	if prefix :=strings.Index(string(tmpBuffer), HeaderPrefix); prefix >= 0 {
		if prefix > 0 {
			tmpBuffer = tmpBuffer[prefix:]
			prefix = 0
		}
		suffix := strings.Index(string(tmpBuffer), HeaderSuffix)
		if suffix > 0 {
			err := json.Unmarshal(tmpBuffer[prefix + HeaderPrefixLen : suffix], &request)
			if err != nil {
				return nil, request, err
			}
			tmpBuffer = tmpBuffer[suffix + HeaderSuffixLen:]
			return tmpBuffer, request, nil
		}
	}
	return tmpBuffer, request, errors.New("Parse Header Error !")
}


func (server *Server) Response (conn net.Conn, status uint16, body string){
	response := Response{status, body}
	responseResult, err := json.Marshal(response)
	responseResult = append(responseResult, ResponseEOF)
	if err !=nil{
		_logger.Warning("Response Result Json Error : ", body)
	}
	_, err = conn.Write(responseResult)
	if err != nil {
		_logger.Warning(err.Error())
	}

}

func (server *Server) Work (conn net.Conn, workChannel chan []byte) (error){

	var responseStatus uint16
	var responseBody string

	defer func(){
		if err :=recover(); err !=nil {
			_logger.Warning("Work Error: %s", err)
		}
	}()
	responseStatus = uint16(StatusOK)
	select {
	case work := <- workChannel:
		var pack Package
		err := json.Unmarshal(work, &pack)
		if err != nil && pack.Request.ContentLength > 0 {
			responseBody = server.Ac.AcFind(pack.Content)
		}else{
			responseStatus = uint16(StatusRequestError)
		}
		server.Response(conn, responseStatus, responseBody)
	}
	return errors.New("Request Empty!" )
}


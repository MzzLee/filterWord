package libs

import (
	"os"
	"bufio"
	"strings"
	"io"
	"fmt"
)

type Node struct {
	Fail *Node
	Child map[uint16]*Node
	Count int
}

var (
	REPLACE_WORD	   = []byte("*")
)

func (node *Node) Insert (keyword string){
	var p *Node
	p = node
	buffer :=[]byte(keyword)
	for _,value := range buffer{
		if p.Child[uint16(value)] == nil{
			p.Child[uint16(value)] = CreateNewNode()
		}
		p = p.Child[uint16(value)]

	}
	p.Count++
}

func (node *Node) ReadLine (fileName string) error {
	fmt.Println("Loading word from :", fileName)
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		if len(line) > 0{
			node.Insert(line)
		}
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
	return nil
}

func (node *Node) BuildAcAutomation (){
	var (
		head int
		tail int
	)
	fmt.Println("Build ac automation ......")
	var q  = make(map[int]*Node)
	node.Fail = nil
	head++
	q[head] = node
	for {
		if head == tail{
			break
		}
		tail++
		temp := q[tail]
		p := new(Node)
		p = nil

		for key,child := range temp.Child{

			if child != nil{
				if temp == node{
					child.Fail = node
				}else{
					p = temp.Fail
					for {
						if p == nil{
							break
						}
						if p.Child[key] != nil {
							child.Fail = p.Child[key]
							break
						}
						p = p.Fail
					}
					if p == nil {
						child.Fail = node
					}
				}
				head++
				q[head] = child
			}
		}
	}

}

func (node *Node) AcFind (buffer []byte) string {
	p := node
	temp := new(Node)
	var start int = 0
	for index, value :=range buffer {
		if IsIgnoreWord(uint16(value)){
			continue
		}
		for{
			if p.Child[uint16(value)] == nil && p != node {
				p = p.Fail
				start = index
			}else{
				break
			}
		}
		p = p.Child[uint16(value)]
		if p == nil{
			p = node
			start = index + 1
		}
		temp = p
		for {
			if temp == node{
				break
			}
			if temp.Count >0{
				for i:=start;i<=index;i++{
					buffer[i] = REPLACE_WORD[0]
				}
			}else{
				break
			}
			temp = temp.Fail
		}
	}
	return string(buffer)
}

func IsIgnoreWord(byte uint16) bool{
	if byte < 48{
		return true
	}
	if byte > 57 && byte < 65{
		return true
	}
	if byte > 90 && byte < 97{
		return true
	}
	if byte > 122 && byte < 128{
		return true
	}
	return false
}

func CreateNewNode() *Node{
	return &Node{nil,make(map[uint16]*Node), 0}
}

func AcBuild (file string) *Node {
	ac := CreateNewNode()
	ac.ReadLine(file)
	ac.BuildAcAutomation()
	return ac
}





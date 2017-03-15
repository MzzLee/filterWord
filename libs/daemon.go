package libs

import (
	"os"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strconv"
)

type Daemon struct{
	Config *Conf
}

func (daemon *Daemon)Run() bool{
	if os.Getppid() != 1 {

		if len(os.Args) < 2 {
			fmt.Println("Place input action For example : start|stop|restart")
			return false
		}
		pid, err := ioutil.ReadFile(daemon.Config.Pid)
		switch os.Args[1] {
		case "start":
			if string(pid) != "" {
				fmt.Println("Server already start !")
				return false
			}

		case "stop":
			if err != nil || string(pid) == "" {
				fmt.Println("Server not start !")
				return false
			}
			os.Remove(daemon.Config.Pid)
			cmd := exec.Command("kill", "-9" , string(pid))
			cmd.Start()
			fmt.Println("Server stop!")
			return false

		case "restart":
			cmd := exec.Command("kill", "-9" , string(pid))
			cmd.Start()
		default:
			fmt.Println("Place input action For example : start|stop|restart")
			return false
		}


		filePath, _ := filepath.Abs(os.Args[0])

		cmd := exec.Command(filePath, os.Args[1], "-c", daemon.Config.File, "-key", daemon.Config.Keyword)
		cmd.Stdin  = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Start()
		return false
	}

	err := ioutil.WriteFile(daemon.Config.Pid, []byte(strconv.Itoa(os.Getpid())), 0666)
	if err != nil{
		fmt.Println("Open Pid Error : ", err.Error())
		return false
	}
	return true
}
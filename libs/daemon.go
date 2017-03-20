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
			fmt.Println("------------------------------------------------------------------")
			fmt.Println("-- Please Input Parameters For Example : start | stop | restart --")
			fmt.Println("------------------------------------------------------------------")
			return false
		}
		pid, err := ioutil.ReadFile(daemon.Config.Pid)
		switch os.Args[1] {
		case "start":
			if string(pid) != "" {
				fmt.Println("------------------------------------------------------------------")
				fmt.Println("------------         Server already start !           ------------")
				fmt.Println("------------------------------------------------------------------")
				return false
			}

		case "stop":
			if err != nil || string(pid) == "" {
				fmt.Println("------------------------------------------------------------------")
				fmt.Println("---------------        Server not start !       ------------------")
				fmt.Println("------------------------------------------------------------------")
				return false
			}
			os.Remove(daemon.Config.Pid)
			cmd := exec.Command("kill", "-9" , string(pid))
			cmd.Start()
			fmt.Println("------------------------------------------------------------------")
			fmt.Println("--------------------       Server stop !        ------------------")
			fmt.Println("------------------------------------------------------------------")
			return false

		case "restart":
			cmd := exec.Command("kill", "-9" , string(pid))
			cmd.Start()
		default:
			fmt.Println("------------------------------------------------------------------")
			fmt.Println("-- Please Input Parameters For Example : start | stop | restart --")
			fmt.Println("------------------------------------------------------------------")
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
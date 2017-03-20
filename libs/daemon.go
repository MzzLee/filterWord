package libs

import (
	"os"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strconv"
	"os/signal"
	"syscall"
)

type Daemon struct{
	Config *Conf
}

func DeamonClient() *Daemon {
	daemon := new(Daemon)
	daemon.Config = GetConfigInstance()
	return daemon
}

func (daemon *Daemon) Run() bool{

	if os.Getppid() != 1 {

		pid, err := ioutil.ReadFile(daemon.Config.Pid)
		switch daemon.Config.Signal {
		case "start":
			if string(pid) != "" {
				fmt.Println("--- Server already start ! ---")
				return false
			}

		case "stop":
			if err != nil || string(pid) == "" {

				fmt.Println("--- Server not start ! ---")

				return false
			}
			os.Remove(daemon.Config.Pid)
			cmd := exec.Command("kill", "-9" , string(pid))
			cmd.Start()
			fmt.Println("--- Server stop ! ---")
			return false

		case "restart":
			cmd := exec.Command("kill", "-9" , string(pid))
			cmd.Start()
		default:

			fmt.Println("--- Please Input -s start | stop | restart ---")
			return false
		}


		filePath, _ := filepath.Abs(os.Args[0])

		cmd := exec.Command(filePath, "-s",   daemon.Config.Signal,
			"-c",   daemon.Config.File,
			"-key", daemon.Config.Keyword,
			"-log", daemon.Config.LogFile)
		cmd.Stdin  = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Start()
		return false
	}
	logger, _ := InitLogger(daemon.Config.LogFile, daemon.Config.Env)
	err := ioutil.WriteFile(daemon.Config.Pid, []byte(strconv.Itoa(os.Getpid())), 0666)
	if err != nil{
		logger.Fatal(err.Error())
		return false
	}
	go daemon.signalListen()
	return true
}

func (daemon *Daemon) signalListen() {

	c :=make(chan os.Signal)
	signal.Notify(c)
	select {
	case s := <-c:
		if s == syscall.SIGHUP{
			os.Remove(daemon.Config.Pid)
			os.Exit(1)
		}
		if s == syscall.SIGQUIT{
			os.Remove(daemon.Config.Pid)
			os.Exit(1)
		}
		if s == syscall.SIGTERM{
			os.Remove(daemon.Config.Pid)
			os.Exit(1)
		}

	}
}

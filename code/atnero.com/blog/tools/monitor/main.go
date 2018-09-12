package main

import (
	"flag"
	"fmt"
	"github.com/hprose/hprose-golang/rpc"
)

type MonitorServices struct {
	EnableBgManager func() string
	StopService     func()
}

var monitorServices *MonitorServices

type FlagArgs struct {
	port string
	cmd  string
}

var flagArgs FlagArgs

func init() {
	flag.StringVar(&flagArgs.port, "port", "8000", "tcp port")
	flag.StringVar(&flagArgs.cmd, "cmd", "", "command")
}

func main() {
	flag.Parse()

	monitorAdd := fmt.Sprintf("tcp4://127.0.0.1:%s/", flagArgs.port)
	fmt.Println("Monitor Address:", monitorAdd)
	client := rpc.NewTCPClient(monitorAdd)
	client.UseService(&monitorServices)
	fmt.Println("Rpc init done")

	switch flagArgs.cmd {
	case "EnableBgManager":
		code := monitorServices.EnableBgManager()
		fmt.Println("Verify code: ", code)
		return
	case "StopService":
		var confirm string
		fmt.Printf("Are your sure to top service (yes/no):")
		fmt.Scanf("%s", &confirm)
		if confirm == "yes" {
			monitorServices.StopService()
			fmt.Println("Stop service command executed success")
			return
		}
		fmt.Println("Quit operation")
		return
	default:
		flag.Usage()
	}
}

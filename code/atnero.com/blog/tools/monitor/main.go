package main

import (
	"flag"
	"fmt"
	"github.com/hprose/hprose-golang/rpc"
)

const (
	ReleaseVersion = "0.1.0"
)

type MonitorServices struct {
	EnableBgManager   func() string
	StopService       func()
	GenInvitationCode func(expireHours int) string
}

var monitorServices *MonitorServices

type FlagArgs struct {
	port  string
	cmd   string
	hours int
}

var flagArgs FlagArgs

func init() {
	flag.StringVar(&flagArgs.port, "port", "8000", "tcp port")
	flag.StringVar(&flagArgs.cmd, "cmd", "", "command: EnableBgManager/StopService/GenInvitationCode")
	flag.IntVar(&flagArgs.hours, "hours", 72, "hours")
}

func main() {
	fmt.Printf("Release Version: %s\n", ReleaseVersion)
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
	case "GenInvitationCode":
		fmt.Printf("Generate invitation code, expire hours = %d\n", flagArgs.hours)
		code := monitorServices.GenInvitationCode(flagArgs.hours)
		fmt.Println(code)
		return
	default:
		flag.Usage()
	}
}

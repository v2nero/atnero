package monitor

import (
	"atnero.com/blog/models/db"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/hprose/hprose-golang/rpc"
)

type MonitorServices struct {
}

func (*MonitorServices) EnableBgManager() string {
	logs.Warning("[Monitor] Eanble background manager")
	return db.DbMgrInst().EnableBgManager()
}
func (*MonitorServices) StopService() {
	logs.Warning("[Monitor] Stop service")
	beego.BeeApp.Server.Close()
}

var tcpServer *rpc.TCPServer

func InitMonitorRpcService() {
	port := beego.AppConfig.String("monitorport")
	monitorAdd := fmt.Sprintf("tcp4://127.0.0.1:%s/", port)
	tcpServer = rpc.NewTCPServer(monitorAdd)
	tcpServer.AddAllMethods(&MonitorServices{})
	go func() {
		logs.Info("[Monitor] address:", monitorAdd)
		err := tcpServer.Start()
		if err != nil {
			logs.Error("[Monitor]", err)
			logs.Error("Quit entire service")
			if beego.BeeApp != nil {
				beego.BeeApp.Server.Close()
			}
		}
	}()
}

func StopService() {
	if tcpServer != nil {
		tcpServer.Stop()
	}
}

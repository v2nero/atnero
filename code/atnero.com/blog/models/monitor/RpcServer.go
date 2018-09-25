package monitor

import (
	"atnero.com/blog/models/db"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/hprose/hprose-golang/rpc"
	"regexp"
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

func (*MonitorServices) GenInvitationCode(expireHours int) string {
	logs.Warning("[Monitor] Generate invitation code")
	return db.DbMgrInst().GenInvitationCode(expireHours)
}

var tcpServer *rpc.TCPServer

type monitorAcceptEvent struct {
}

func (this *monitorAcceptEvent) OnAccept(context *rpc.SocketContext) error {
	strAddr := context.RemoteAddr().String()
	logs.Info("[Monitor] address ", strAddr)
	matched, err := regexp.MatchString(`^::\d{1,3}:\d*`, strAddr)
	if err != nil {
		return err
	}
	if matched {
		return nil
	}
	matched, err = regexp.MatchString(`^127.0.0.\d{1,3}:\d*`, strAddr)
	if err != nil {
		return err
	}
	if matched {
		return nil
	}
	return fmt.Errorf("remote address must be local ip")
}

func InitMonitorRpcService() {
	port := beego.AppConfig.String("monitorport")
	monitorAdd := fmt.Sprintf("tcp4://127.0.0.1:%s/", port)
	tcpServer = rpc.NewTCPServer(monitorAdd)
	tcpServer.AddAllMethods(&MonitorServices{})
	tcpServer.Event = &monitorAcceptEvent{}
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

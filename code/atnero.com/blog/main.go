package main

import (
	_ "atnero.com/blog/models"
	_ "atnero.com/blog/models/db"
	"atnero.com/blog/models/monitor"
	_ "atnero.com/blog/routers"
	"github.com/astaxie/beego"
	_ "github.com/astaxie/beego/logs"
)

func main() {
	//写用户权限管理
	monitor.InitMonitorRpcService()
	beego.AddTemplateExt("html")
	beego.Run()
}

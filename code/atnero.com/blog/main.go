package main

import (
	_ "atnero.com/blog/models"
	_ "atnero.com/blog/models/db"
	"atnero.com/blog/models/monitor"
	_ "atnero.com/blog/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

func main() {
	//TODO: 写monitor客户端
	if err := monitor.InitServer(); err != nil {
		logs.Error(err)
		return
	}
	beego.AddTemplateExt("html")
	beego.Run()
}

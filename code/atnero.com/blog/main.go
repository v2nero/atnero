package main

import (
	"atnero.com/blog/models"
	_ "atnero.com/blog/models/db"
	"atnero.com/blog/models/monitor"
	_ "atnero.com/blog/routers"
	"github.com/astaxie/beego"
	_ "github.com/astaxie/beego/logs"
)

func main() {
	//ckeditor4
	ckeditor4Path := beego.AppConfig.String("ckeditor4")
	if len(ckeditor4Path) == 0 {
		beego.Error("[CKEDITOR4] missing ckeditor4")
		return
	}
	beego.Info("[CKEDITOR4] path=", ckeditor4Path)
	beego.SetStaticPath("/thirdparty/ckeditor4", ckeditor4Path)

	models.CheckRightSetDependencies()
	monitor.InitMonitorRpcService()
	beego.AddTemplateExt("html")
	beego.AddTemplateExt("xml")

	beego.Run()

	//TODO:
	//1. 分类文章列表
	//2. 图片上传
}

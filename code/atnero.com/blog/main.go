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

	beego.Run()

	//TODO:
	//1. 文章分类增加
	//2. 暂不做标签, 后续修改
	//3. 文章创建，修改,删除
}

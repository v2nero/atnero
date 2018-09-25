package main

import (
	"atnero.com/blog/models"
	_ "atnero.com/blog/models/db"
	"atnero.com/blog/models/monitor"
	_ "atnero.com/blog/routers"
	"github.com/astaxie/beego"
	_ "github.com/astaxie/beego/logs"
)

const (
	ReleaseVersion = "0.1.1"
)

func main() {
	beego.Info("[Version] ReleaseVersion ", ReleaseVersion)
	//ckeditor4
	ckeditor4Path := beego.AppConfig.String("ckeditor4")
	if len(ckeditor4Path) == 0 {
		beego.Error("[CKEDITOR4] missing ckeditor4")
		return
	}
	beego.Info("[CKEDITOR4] path=", ckeditor4Path)
	beego.SetStaticPath("/thirdparty/ckeditor4", ckeditor4Path)

	ck_imgupload_viewdir := beego.AppConfig.String("ckeditor::image_upload_view_dir")
	ck_imgupload_dir := beego.AppConfig.String("ckeditor::image_upload_dir")
	if len(ck_imgupload_viewdir) == 0 {
		beego.Error("[CKEDITOR4] missing ckeditor::image_upload_view_dir")
		return
	}
	if len(ck_imgupload_dir) == 0 {
		beego.Error("[CKEDITOR4] missing ckeditor::image_upload_dir")
		return
	}

	beego.SetStaticPath(ck_imgupload_viewdir, ck_imgupload_dir)

	models.CheckRightSetDependencies()
	monitor.InitMonitorRpcService()
	beego.AddTemplateExt("html")
	beego.AddTemplateExt("xml")

	beego.Run()
}

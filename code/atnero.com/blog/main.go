package main

import (
	_ "atnero.com/blog/db"
	_ "atnero.com/blog/routers"
	"github.com/astaxie/beego"
)

func main() {
	//TODO:
	//整理代码管理
	//整理代码命名

	beego.AddTemplateExt("html")
	beego.Run()
}

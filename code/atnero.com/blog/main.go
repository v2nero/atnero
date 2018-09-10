package main

import (
	_ "atnero.com/blog/models"
	_ "atnero.com/blog/models/db"
	_ "atnero.com/blog/routers"
	"github.com/astaxie/beego"
)

func main() {
	beego.AddTemplateExt("html")
	beego.Run()
}

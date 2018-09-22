package controllers

import (
	_ "github.com/astaxie/beego"
)

type MainController struct {
	CommonPageContainerController
}

func (this *MainController) Get() {

	this.InitLayout()
	this.TplName = "index.tpl"
	this.Data["Title"] = "@Nero"
	this.Data["Website"] = "atnero.com"
	this.Data["Email"] = ""
}

package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "atnero.com"
	c.Data["Email"] = ""
	c.TplName = "index.tpl"
}

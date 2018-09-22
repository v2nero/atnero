package controllers

import (
	//"github.com/astaxie/beego"
	blogSess "atnero.com/blog/models/session"
)

type LogoutController struct {
	CommonPageContainerController
}

func (this *LogoutController) Get() {
	this.InitLayout()
	this.TplName = "login/logout.html"
	this.Data["Title"] = "登出"
	blogSess.Logout(&this.Controller)
	this.LayoutSections["Scripts"] = "common/Redirection.html"
	this.Data["RedirectionURL"] = "/"
	return
}

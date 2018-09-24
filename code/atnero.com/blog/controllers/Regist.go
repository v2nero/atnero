package controllers

import (
	blogSess "atnero.com/blog/models/session"
)

type RegistController struct {
	CommonPageContainerController
}

func (this *RegistController) Get() {
	this.InitLayout()
	this.TplName = "login/regist.html"
	this.Data["PageTitle"] = "注册"
	if blogSess.Logined(&this.Controller) {
		this.Data["AlreadyLogin"] = true
		return
	}
}

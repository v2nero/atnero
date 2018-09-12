package controllers

import (
	_ "atnero.com/blog/models"
	_ "atnero.com/blog/models/db"
	blogSess "atnero.com/blog/models/session"
)

type ManagerUserRightController struct {
	CommonPageContainerController
}

func (this *ManagerUserRightController) Get() {
	if !blogSess.BgManagerEnabled(&this.Controller) {
		this.Abort("404")
		return
	}
	this.InitLayout()
	this.TplName = "manager/userrights.html"
	this.Data["Title"] = "[后台管理]用户权限 @Nero"
}
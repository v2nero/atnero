package controllers

import (
	//"github.com/astaxie/beego"
)

type ManagerUserRightController struct {
	CommonPageContainerController
}

func (this *ManagerUserRightController) Get() {
	this.InitLayout()
	this.TplName = "manager/userright.html"
	this.Data["Title"] = "[后台管理]用户权限 @Nero"
}
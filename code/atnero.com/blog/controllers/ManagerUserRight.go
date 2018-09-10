package controllers

import (
	_ "atnero.com/blog/models"
	"atnero.com/blog/models/db"
)

type ManagerUserRightController struct {
	CommonPageContainerController
}

func (this *ManagerUserRightController) Get() {
	if !db.DbMgrInst().GetBgManagerEnable() {
		this.Abort("404")
		return
	}
	this.InitLayout()
	this.TplName = "manager/userrights.html"
	this.Data["Title"] = "[后台管理]用户权限 @Nero"
}
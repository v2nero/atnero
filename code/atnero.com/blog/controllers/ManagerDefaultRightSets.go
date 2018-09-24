package controllers

import (
	"atnero.com/blog/models"
	_ "atnero.com/blog/models/db"
	blogSess "atnero.com/blog/models/session"
)

type ManagerDefaultUserRightSetsController struct {
	CommonPageContainerController
}

func (this *ManagerDefaultUserRightSetsController) Get() {
	if !models.TestManagerInst().BgManagerTestEnabled(&this.Controller) {
		if !blogSess.BgManagerEnabled(&this.Controller) {
			this.Abort("404")
			return
		}
	}
	this.InitLayout()
	this.TplName = "manager/defaultuserrightsets.html"
	this.Data["PageTitle"] = "[后台管理]默认权限配置 @Nero"

	mng := models.UserRightsMngInst()
	exposeDefaultRightSets := mng.GetDefaultRightSetsList()
	exposeRightSets := mng.GetRightSets()
	this.Data["DefaultRightSets"] = exposeDefaultRightSets
	this.Data["TotalRightSets"] = exposeRightSets
}

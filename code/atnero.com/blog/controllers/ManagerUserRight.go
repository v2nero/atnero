package controllers

import (
	"atnero.com/blog/models"
	_ "atnero.com/blog/models/db"
	blogSess "atnero.com/blog/models/session"
)

type ManagerUserRightController struct {
	CommonPageContainerController
}

type userRightItemExpose struct {
	Name    string
	Enabled bool
	Dsc     string
}

func (this *ManagerUserRightController) Get() {
	if !models.TestManagerInst().BgManagerTestEnabled(&this.Controller) {
		if !blogSess.BgManagerEnabled(&this.Controller) {
			this.Abort("404")
			return
		}
	}
	this.InitLayout()
	this.TplName = "manager/userright.html"
	this.Data["Title"] = "[后台管理]用户权限 @Nero"

	//权限项
	items := []userRightItemExpose{}
	rightItems := models.UserRightsMngInst().GetRightItems()
	for _, v := range rightItems {
		enabled, err := models.UserRightsMngInst().RightItemEnabled(v)
		if err != nil {
			continue
		}
		dsc, bExist := models.UserRightsMngInst().GetRightItemDiscription(v)
		if !bExist {
			continue
		}
		items = append(items, userRightItemExpose{v, enabled, dsc})
	}
	this.Data["RightItems"] = items
}
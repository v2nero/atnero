package controllers

import (
	"atnero.com/blog/models"
	_ "atnero.com/blog/models/db"
	blogSess "atnero.com/blog/models/session"
)

type ManagerUserRightSetsController struct {
	CommonPageContainerController
}

type userRightSetExpose struct {
	Name       string
	Dsc        string
	RightItems []string
}

func (this *ManagerUserRightSetsController) Get() {
	if !models.TestManagerInst().BgManagerTestEnabled(&this.Controller) {
		if !blogSess.BgManagerEnabled(&this.Controller) {
			this.Abort("404")
			return
		}
	}
	this.InitLayout()
	this.TplName = "manager/userrightsets.html"
	this.Data["Title"] = "[后台管理]用户权限集 @Nero"

	//权限项
	exposeSets := []userRightSetExpose{}
	rightSets := models.UserRightsMngInst().GetRightSets()
	for _, v := range rightSets {
		dsc, bExist := models.UserRightsMngInst().GetRightSetDiscription(v)
		if !bExist {
			continue
		}
		rightItems := []string{}
		for _, v := range models.UserRightsMngInst().GetRightSetRightItems(v) {
			rightItems = append(rightItems, v)
		}
		exposeSets = append(exposeSets, userRightSetExpose{v, dsc, rightItems})
	}
	this.Data["RightSets"] = exposeSets
}

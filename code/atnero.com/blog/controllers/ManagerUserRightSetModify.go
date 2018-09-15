package controllers

import (
	"atnero.com/blog/models"
	_ "atnero.com/blog/models/db"
	blogSess "atnero.com/blog/models/session"
)

type ManagerUserRightSetModifyController struct {
	CommonPageContainerController
}

func (this *ManagerUserRightSetModifyController) Get() {
	if !models.TestManagerInst().BgManagerTestEnabled(&this.Controller) {
		if !blogSess.BgManagerEnabled(&this.Controller) {
			this.Abort("404")
			return
		}
	}
	this.InitLayout()
	this.TplName = "manager/userrightset_modify.html"
	this.Data["Title"] = "[后台管理]用户权限集修改 @Nero"
	rightMng := models.UserRightsMngInst()

	setName := this.Ctx.Input.Param(":setname")
	if !rightMng.HasRightSet(setName) {
		this.Data["HasSuchSet"] = false
		return
	}
	this.Data["HasSuchSet"] = true

	this.Data["Name"] = setName
	setDsc, _ := rightMng.GetRightSetDiscription(setName)
	this.Data["Dsc"] = setDsc

	hadRightItems := rightMng.GetRightSetRightItems(setName)

	allRightItemsExpose := []userRightItemExpose{}
	allRightItems := models.UserRightsMngInst().GetRightItems()
	for _, v := range allRightItems {
		dsc, bExist := rightMng.GetRightItemDiscription(v)
		if !bExist {
			continue
		}
		bHasIt := false
		for _, iname := range hadRightItems {
			if v == iname {
				bHasIt = true
				break
			}
		}
		allRightItemsExpose = append(allRightItemsExpose, userRightItemExpose{v, bHasIt, dsc})
	}
	this.Data["RightItems"] = allRightItemsExpose
}

package controllers

import (
	"atnero.com/blog/models"
	blogSess "atnero.com/blog/models/session"
	"github.com/astaxie/beego"
	//"strconv"
)

type ApiDefaultRightSetController struct {
	beego.Controller
}

func (this *ApiDefaultRightSetController) handleUpdate() bool {
	strName := this.GetString("name")
	if len(strName) == 0 {
		return false
	}
	strDsc := this.GetString("dsc")
	strSetName := this.GetString("setname")
	if len(strSetName) == 0 {
		return false
	}
	err := models.UserRightsMngInst().UpdateDefaultRightSet(
		strName, strDsc, strSetName)
	if err != nil {
		return false
	}
	return true
}

func (this *ApiDefaultRightSetController) handleCreate() bool {
	strName := this.GetString("name")
	if len(strName) == 0 {
		return false
	}
	strDsc := this.GetString("dsc")
	strSetName := this.GetString("setname")
	if len(strSetName) == 0 {
		return false
	}
	err := models.UserRightsMngInst().AddDefaultRightSet(
		strName, strDsc, strSetName)
	if err != nil {
		return false
	}
	return true
}

func (this *ApiDefaultRightSetController) Post() {
	if !models.TestManagerInst().BgManagerTestEnabled(&this.Controller) {
		if !blogSess.BgManagerEnabled(&this.Controller) {
			this.Abort("404")
			return
		}
	}

	result := false

	cmd := this.GetString("cmd")
	switch cmd {
	case "update":
		result = this.handleUpdate()
		break
	case "create":
		result = this.handleCreate()
		break
	}

	if result {
		this.Data["result"] = "success"
	} else {
		this.Data["result"] = "fail"
	}
	this.TplName = "manager/ApiDefaultRightSet.html"
}

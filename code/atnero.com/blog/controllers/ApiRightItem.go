package controllers

import (
	"atnero.com/blog/models"
	blogSess "atnero.com/blog/models/session"
	"github.com/astaxie/beego"
	//"strconv"
)

type ApiRightItemController struct {
	beego.Controller
}

/*
type formRightItemUpdateInfo struct {
	Type   string `form:"type"`
	Id     int64  `form:"id"`
	Name   string `form:"name"`
	Enabled bool   `form:"enabled"`
	Dsc    string `form:"dsc"`
}
*/

func handleEnable(this *ApiRightItemController) bool {
	strName := this.GetString("name")
	if len(strName) == 0 {
		return false
	}
	bEnabled := false
	strEnabled := this.GetString("enabled")
	if strEnabled == "true" {
		bEnabled = true
	} else if strEnabled == "false" {
		bEnabled = false
	} else {
		return false
	}
	err := models.UserRightsMngInst().EnableRightItem(strName, bEnabled)
	if err != nil {
		return false
	}
	return true
}

func handleCreate(this *ApiRightItemController) bool {
	strName := this.GetString("name")
	if len(strName) == 0 {
		return false
	}
	bEnabled := false
	strEnabled := this.GetString("enabled")
	if strEnabled == "true" {
		bEnabled = true
	} else if strEnabled == "false" {
		bEnabled = false
	} else {
		return false
	}
	strDsc := this.GetString("dsc")
	err := models.UserRightsMngInst().AddRightItem(strName, bEnabled, strDsc)
	if err != nil {
		return false
	}
	return true
}

func (this *ApiRightItemController) Post() {
	if !blogSess.BgManagerEnabled(&this.Controller) {
		this.Abort("404")
		return
	}

	result := false

	cmd := this.GetString("cmd")
	switch cmd {
	case "enable":
		result = handleEnable(this)
		break
	case "modify":
		break
	case "create":
		result = handleCreate(this)
		break
	}

	if result {
		this.Data["result"] = "success"
	} else {
		this.Data["result"] = "fail"
	}
	this.TplName = "manager/ApiRightItem.html"
}

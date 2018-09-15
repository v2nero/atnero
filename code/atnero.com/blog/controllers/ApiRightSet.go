package controllers

import (
	"atnero.com/blog/models"
	blogSess "atnero.com/blog/models/session"
	"github.com/astaxie/beego"
	//"strconv"
)

type ApiRightSetController struct {
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

func (this *ApiRightSetController) handleEnable() bool {
	strName := this.GetString("name")
	if len(strName) == 0 {
		return false
	}
	strItemName := this.GetString("rightitemname")
	if len(strItemName) == 0 {
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
	var err error
	if bEnabled {
		err = models.UserRightsMngInst().AddRightItem2RightSet(strName, strItemName)
	} else {
		err = models.UserRightsMngInst().DelRightItemFromRightSet(strName, strItemName)
	}
	if err != nil {
		return false
	}
	return true
}

func (this *ApiRightSetController) handleCreate() bool {
	strName := this.GetString("name")
	if len(strName) == 0 {
		return false
	}
	strDsc := this.GetString("dsc")
	err := models.UserRightsMngInst().AddRightSet(strName, strDsc)
	if err != nil {
		return false
	}
	return true
}

func (this *ApiRightSetController) Post() {
	if !models.TestManagerInst().BgManagerTestEnabled(&this.Controller) {
		if !blogSess.BgManagerEnabled(&this.Controller) {
			this.Abort("404")
			return
		}
	}

	result := false

	cmd := this.GetString("cmd")
	switch cmd {
	case "enable":
		result = this.handleEnable()
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
	this.TplName = "manager/ApiRightSet.html"
}

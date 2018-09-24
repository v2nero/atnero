package controllers

import (
	"atnero.com/blog/models/db"
	blogSess "atnero.com/blog/models/session"
)

type ManagerVerifyController struct {
	CommonPageContainerController
}

func (this *ManagerVerifyController) Get() {
	this.InitLayout()
	this.TplName = "manager/verify.html"
	this.Data["Title"] = "[后台管理]验证"
	if blogSess.BgManagerEnabled(&this.Controller) {
		this.Data["AlreadyVerified"] = true
		return
	}

	if !db.DbMgrInst().GetBgManagerEnable() {
		this.Abort("404")
		return
	}
}

type verifyInfo struct {
	Code string `form:"verifycode"`
}

func (this *ManagerVerifyController) Post() {
	this.InitLayout()
	this.TplName = "manager/verify.html"
	this.Data["PageTitle"] = "[后台管理]验证"
	if blogSess.BgManagerEnabled(&this.Controller) {
		this.Data["AlreadyVerified"] = true
		return
	}
	if !db.DbMgrInst().GetBgManagerEnable() {
		this.Abort("404")
		return
	}
	verifyInfo := verifyInfo{}
	if err := this.ParseForm(&verifyInfo); err != nil {
		this.Data["InputDataError"] = true
		return
	}
	if blogSess.EnableBgManager(&this.Controller, verifyInfo.Code) {
		this.Data["AlreadyVerified"] = true
		return
	} else {
		this.Data["VerifyFail"] = true
		return
	}
}

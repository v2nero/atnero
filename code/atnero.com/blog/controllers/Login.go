package controllers

import (
	//"github.com/astaxie/beego"
	blogSess "atnero.com/blog/models/session"
)

type LoginController struct {
	CommonPageContainerController
}

func (this *LoginController) Get() {
	this.InitLayout()
	this.TplName = "login/login.html"
	this.Data["PageTitle"] = "登陆"
	if blogSess.Logined(&this.Controller) {
		this.Data["AlreadyLogin"] = true
		return
	}
}

type loginInfo struct {
	Name string `form:"user"`
	Pwd  string `form:"password"`
}

func (this *LoginController) Post() {
	this.InitLayout()
	this.TplName = "login/login.html"
	this.Data["PageTitle"] = "登陆"

	//已经登陆，显示错误
	if blogSess.Logined(&this.Controller) {
		this.Data["AlreadyLogin"] = true
		return
	}

	//读取登陆信息
	loginInfo := loginInfo{}
	if err := this.ParseForm(&loginInfo); err != nil {
		this.Data["InputDataError"] = true
		return
	}
	//用户名不正常
	if !checkUserName(loginInfo.Name) {
		this.Data["InputDataError"] = true
		return
	}
	if !checkPwd(loginInfo.Pwd) {
		this.Data["InputDataError"] = true
		return
	}

	//设置session
	if blogSess.Login(&this.Controller, loginInfo.Name, loginInfo.Pwd) {
		this.LayoutSections["Scripts"] = "common/Redirection.html"
		this.Data["RedirectionURL"] = "/"
		return
	} else {
		this.Data["LoginFail"] = true
		this.Data["LoginFailTime"] = blogSess.GetLoginFailInterval()
	}
}

package controllers

import (
	//"github.com/astaxie/beego"
)

type LoginController struct {
	CommonPageContainerController
}

func (this *LoginController) Get() {
	this.InitLayout()
	this.TplName = "login/login.html"
	this.Data["Title"] = "登陆 @Nero"
	userInfo := this.GetSession("user")
	if userInfo == nil {
		return
	}
	this.Data["AlreadyLogin"] = true
}

type loginInfo struct {
	Name string `form:"user"`
	Pwd string `form:"password"`
}

func (this *LoginController) Post() {
	this.InitLayout()
	this.TplName = "login/login.html"
	this.Data["Title"] = "登陆 @Nero"

	//已经登陆，显示错误
	userInfo := this.GetSession("user")
	if userInfo != nil {
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
	if len(loginInfo.Name)<4 || len(loginInfo.Name)>24 {
		this.Data["InputDataError"] = true;
		return
	}
	
	//设置session
	sessUserInfo := make(map[string]interface{})
	sessUserInfo["name"] = loginInfo.Name
	this.SetSession("user", sessUserInfo)
	this.LayoutSections["Scripts"] = "common/Redirection.html"
	this.Data["RedirectionURL"] = "/"
}

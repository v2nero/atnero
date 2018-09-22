package controllers

import (
	_ "atnero.com/blog/models"
	_ "atnero.com/blog/models/db"
	blogSess "atnero.com/blog/models/session"
)

type ArticleCreateController struct {
	CommonPageContainerController
}

func (this *ArticleCreateController) Get() {
	if !blogSess.Logined(&this.Controller) {
		this.Abort("404")
		return
	}
	//TODO: 增加权限控制
	this.InitLayout()
	this.TplName = "article/create.html"
	this.Data["Title"] = "博客创建 @Nero"
}

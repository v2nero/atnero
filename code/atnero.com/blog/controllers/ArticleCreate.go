package controllers

import (
	"atnero.com/blog/models"
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

	if !blogSess.UserHasRight(&this.Controller, "create_article") {
		this.Abort("404")
		return
	}

	this.InitLayout()
	this.TplName = "article/create.html"
	this.Data["PageTitle"] = "博客创建 @Nero"
	this.Data["ArticleSorts"] = models.ArticleManagerInst().GetSorts()
	this.Data["ArticleClasses"] = models.ArticleManagerInst().GetClasses()
}

package controllers

import (
	"atnero.com/blog/models"
	_ "atnero.com/blog/models/db"
	blogSess "atnero.com/blog/models/session"
)

type ManagerArticleClassesController struct {
	CommonPageContainerController
}

func (this *ManagerArticleClassesController) Get() {
	if !models.TestManagerInst().BgManagerTestEnabled(&this.Controller) {
		if !blogSess.BgManagerEnabled(&this.Controller) {
			this.Abort("404")
			return
		}
	}
	this.InitLayout()
	this.TplName = "manager/articlesorts.html"
	this.Data["Title"] = "[后台管理]文章归类 @Nero"

	mng := models.ArticleManagerInst()
	this.Data["Items"] = mng.GetClasses()
	this.Data["Explain"] = "文章归类"
	this.Data["FormActionUrl"] = "/manager/api_articleclass"
}

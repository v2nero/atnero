package controllers

import (
	"atnero.com/blog/models"
	blogSess "atnero.com/blog/models/session"
)

type ManagerArticleSortsController struct {
	CommonPageContainerController
}

type articleSortExpose struct {
	Name string
}

func (this *ManagerArticleSortsController) Get() {
	if !models.TestManagerInst().BgManagerTestEnabled(&this.Controller) {
		if !blogSess.BgManagerEnabled(&this.Controller) {
			this.Abort("404")
			return
		}
	}
	this.InitLayout()
	this.TplName = "manager/articlesorts.html"
	this.Data["PageTitle"] = "[后台管理]文章类型 @Nero"

	mng := models.ArticleManagerInst()
	this.Data["Items"] = mng.GetSorts()
	this.Data["Explain"] = "文章类型"
	this.Data["FormActionUrl"] = "/manager/api_articlesort"
}

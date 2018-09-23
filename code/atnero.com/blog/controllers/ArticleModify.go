package controllers

import (
	"atnero.com/blog/models"
	_ "atnero.com/blog/models/db"
	blogSess "atnero.com/blog/models/session"
	"strconv"
)

type ArticleModifyController struct {
	CommonPageContainerController
}

func (this *ArticleModifyController) Get() {
	if !blogSess.Logined(&this.Controller) {
		this.Abort("404")
		return
	}

	strArticleId := this.Ctx.Input.Param(":id")
	articleId, err := strconv.ParseInt(strArticleId, 10, 64)
	if err != nil {
		this.Abort("404")
		return
	}

	articleData, err := models.ArticleManagerInst().GetArticleData(articleId)
	if err != nil {
		this.Abort("404")
		return
	}

	_, userId, err := blogSess.GetUserBaseInfo(&this.Controller)
	if err != nil {
		this.Abort("404")
		return
	}

	if userId != articleData.UserId {
		if !blogSess.UserHasRight(&this.Controller, "edit_others_article") {
			this.Abort("404")
			return
		}
	} else if !blogSess.UserHasRight(&this.Controller, "edit_my_article") {
		this.Abort("404")
		return
	}

	this.InitLayout()
	this.TplName = "article/modify.html"
	this.Data["Title"] = "博客创建 @Nero"
	this.Data["ArticleSorts"] = models.ArticleManagerInst().GetSorts()
	this.Data["ArticleClasses"] = models.ArticleManagerInst().GetClasses()
	this.Data["ArticleDataView"] = articleData
}

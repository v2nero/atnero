package controllers

import (
	"atnero.com/blog/models"
	blogSess "atnero.com/blog/models/session"
	"strconv"
)

type ArticleViewController struct {
	CommonPageContainerController
}

func (this *ArticleViewController) Get() {
	this.InitLayout()
	this.TplName = "article/view.html"
	this.Data["Title"] = "读博文 @Nero"

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

	bLogin := false
	_, userId, err := blogSess.GetUserBaseInfo(&this.Controller)
	if err == nil {
		bLogin = true
	}

	var viewRightItem string
	var editRightItem string
	if !bLogin || userId != articleData.UserId {
		if articleData.Published {
			viewRightItem = "view_others_published_article"
		} else {
			viewRightItem = "view_others_hidden_article"
		}
		editRightItem = "edit_others_article"
	} else {
		if articleData.Published {
			viewRightItem = "view_my_published_article"
		} else {
			viewRightItem = "view_my_hidden_article"
		}
		editRightItem = "edit_my_article"
	}
	if !blogSess.UserHasRight(&this.Controller, viewRightItem) {
		this.Abort("404")
		return
	}

	if blogSess.UserHasRight(&this.Controller, editRightItem) {
		this.Data["Editable"] = true
	}
	this.Data["ArticleDataView"] = articleData
}

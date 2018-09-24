package controllers

import (
	"atnero.com/blog/models"
	blogSess "atnero.com/blog/models/session"
	"strconv"
)

type ArticleListController struct {
	CommonPageContainerController
}

func (this *ArticleListController) Get() {
	this.InitLayout()
	this.TplName = "article/list.html"
	this.Data["PageTitle"] = "博客列表 @Nero"

	if !blogSess.UserHasRight(&this.Controller, "view_others_published_article") {
		this.Data["OutOfRight"] = true
		return
	}

	strPageIndex := this.Ctx.Input.Param(":index")
	pageIndex, err := strconv.ParseInt(strPageIndex, 10, 64)
	if err != nil {
		pageIndex = 0
	}

	var limit int64
	limit = 10

	articleShortViews, err := models.ArticleManagerInst().GetArticlesShortViewOfAll(true, pageIndex, limit)
	if err != nil {
		this.Data["InternalError"] = true
		return
	}

	this.Data["ArticleShortViews"] = articleShortViews
	this.pageList(pageIndex, limit)
}

func (this *ArticleListController) pageList(curIndex int64, limit int64) {
	total, err := models.ArticleManagerInst().GetArticlesNumOfAll(true)
	if err != nil {
		return
	}
	pages := int64(total) / limit
	if int64(total)%limit > 0 {
		pages = pages + 1
	}
	this.Data["TotalPages"] = pages
	this.Data["CurrentPageIndex"] = curIndex
}

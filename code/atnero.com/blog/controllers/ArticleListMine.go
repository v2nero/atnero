package controllers

import (
	"atnero.com/blog/models"
	blogSess "atnero.com/blog/models/session"
	"strconv"
)

type ArticleListMineController struct {
	CommonPageContainerController
}

func (this *ArticleListMineController) Get() {

	if !blogSess.Logined(&this.Controller) {
		this.Abort("404")
		return
	}

	this.InitLayout()
	this.TplName = "article/listmine.html"
	this.Data["PageTitle"] = "博客列表 @Nero"

	_, userId, err := blogSess.GetUserBaseInfo(&this.Controller)
	if err != nil {
		this.Data["InternalError"] = true
		return
	}

	if !blogSess.UserHasRight(&this.Controller, "view_my_published_article") {
		this.Data["OutOfRight"] = true
		return
	}

	viewHidden := blogSess.UserHasRight(&this.Controller, "view_my_hidden_article")

	strPageIndex := this.Ctx.Input.Param(":index")
	pageIndex, err := strconv.ParseInt(strPageIndex, 10, 64)
	if err != nil {
		pageIndex = 0
	}

	var limit int64
	limit = 10

	articleShortViews, err := models.ArticleManagerInst().GetArticlesShortViewOfUser(userId, !viewHidden, pageIndex, limit)
	if err != nil {
		this.Data["InternalError"] = true
		return
	}

	userRightEdit := blogSess.UserHasRight(&this.Controller, "edit_my_article")

	this.Data["ArticleShortViews"] = articleShortViews
	this.Data["Editable"] = userRightEdit
	this.pageList(userId, pageIndex, limit, !viewHidden)
}

func (this *ArticleListMineController) pageList(userId int64, curIndex int64, limit int64, publishedOnly bool) {
	total, err := models.ArticleManagerInst().GetArticlesNumOfUser(userId, publishedOnly)
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

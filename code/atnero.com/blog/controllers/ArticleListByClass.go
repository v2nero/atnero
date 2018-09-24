package controllers

import (
	"atnero.com/blog/models"
	blogSess "atnero.com/blog/models/session"
	"strconv"
)

type ArticleListByClassController struct {
	CommonPageContainerController
}

func (this *ArticleListByClassController) Get() {
	this.InitLayout()
	this.TplName = "article/listbyclass.html"
	this.Data["PageTitle"] = "博客列表 @Nero"

	if !blogSess.UserHasRight(&this.Controller, "view_others_published_article") {
		this.Data["OutOfRight"] = true
		return
	}

	strClassId := this.Ctx.Input.Param(":classid")
	classId, err := strconv.ParseInt(strClassId, 10, 64)
	if err != nil {
		this.Abort("404")
		return
	}
	strClassName, err := models.ArticleManagerInst().GetClassName(classId)
	if err != nil {
		this.Abort("404")
		return
	}

	strPageIndex := this.Ctx.Input.Param(":index")
	pageIndex, err := strconv.ParseInt(strPageIndex, 10, 64)
	if err != nil {
		pageIndex = 0
	}

	var limit int64
	limit = 10

	articleShortViews, err := models.ArticleManagerInst().GetArticlesShortViewOfClass(classId, pageIndex, limit)
	if err != nil {
		this.Data["InternalError"] = true
		return
	}

	this.Data["ArticleShortViews"] = articleShortViews
	this.Data["ClassName"] = strClassName
	this.Data["ClassId"] = classId
	this.pageList(classId, pageIndex, limit)
}

func (this *ArticleListByClassController) pageList(classId int64, curIndex int64, limit int64) {
	total, err := models.ArticleManagerInst().GetArticlesNumOfClass(classId)
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

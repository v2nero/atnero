package controllers

import (
	"atnero.com/blog/models"
	blogSess "atnero.com/blog/models/session"
	"github.com/astaxie/beego"
	"strconv"
)

type ApiArticleController struct {
	beego.Controller
}

func (this *ApiArticleController) handleCreate() bool {
	strTitle := this.GetString("title")
	if len(strTitle) == 0 {
		return false
	}

	strSort := this.GetString("sort")
	if len(strSort) == 0 {
		return false
	}

	strClass := this.GetString("class")
	if len(strClass) == 0 {
		return false
	}

	bPublish := false
	strPublish := this.GetString("publish")
	if strPublish == "true" {
		bPublish = true
	} else if strPublish == "false" {
		bPublish = false
	} else {
		return false
	}
	strContent := this.GetString("content")
	if len(strContent) == 0 {
		return false
	}
	if !blogSess.UserHasRight(&this.Controller, "create_article") {
		return false
	}
	_, userId, err := blogSess.GetUserBaseInfo(&this.Controller)
	if err != nil {
		return false
	}

	articleId, err := models.ArticleManagerInst().AddArticle(
		userId, strTitle, strContent, strSort, strClass, bPublish)
	if err != nil {
		return false
	}
	this.Data["ArticleId"] = articleId
	return true
}

func (this *ApiArticleController) handleModify() bool {
	strId := this.GetString("id")
	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		return false
	}

	strTitle := this.GetString("title")
	if len(strTitle) == 0 {
		return false
	}

	strSort := this.GetString("sort")
	if len(strSort) == 0 {
		return false
	}

	strClass := this.GetString("class")
	if len(strClass) == 0 {
		return false
	}

	bPublish := false
	strPublish := this.GetString("publish")
	if strPublish == "true" {
		bPublish = true
	} else if strPublish == "false" {
		bPublish = false
	} else {
		return false
	}
	strContent := this.GetString("content")
	if len(strContent) == 0 {
		return false
	}

	if !blogSess.UserHasRight(&this.Controller, "create_article") {
		return false
	}

	strRightRequired := "edit_my_article"
	_, userId, err := blogSess.GetUserBaseInfo(&this.Controller)
	if err != nil {
		return false
	}

	oldData, err := models.ArticleManagerInst().GetArticleData(id)
	if err != nil || oldData.UserId != userId {
		strRightRequired = "edit_others_article"
	}
	if !blogSess.UserHasRight(&this.Controller, strRightRequired) {
		return false
	}

	err = models.ArticleManagerInst().UpdateArticle(id, strTitle, strContent, strSort, strClass, bPublish)
	if err != nil {
		return false
	}
	this.Data["ArticleId"] = id
	return true
}

func (this *ApiArticleController) Post() {
	if !models.TestManagerInst().BgManagerTestEnabled(&this.Controller) {
		if !blogSess.BgManagerEnabled(&this.Controller) {
			this.Abort("404")
			return
		}
	}

	result := false

	cmd := this.GetString("cmd")
	switch cmd {
	case "create":
		result = this.handleCreate()
		break
	case "modify":
		result = this.handleModify()
		break
	}

	if result {
		this.Data["result"] = "success"
	} else {
		this.Data["result"] = "fail"
	}
	this.TplName = "article/ApiCreate.xml"
}

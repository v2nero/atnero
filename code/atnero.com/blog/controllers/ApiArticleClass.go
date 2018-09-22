package controllers

import (
	"atnero.com/blog/models"
	blogSess "atnero.com/blog/models/session"
	"github.com/astaxie/beego"
	//"strconv"
)

type ApiArticleClassController struct {
	beego.Controller
}

func (this *ApiArticleClassController) handleAdd() bool {
	strName := this.GetString("name")
	if len(strName) == 0 {
		return false
	}

	err := models.ArticleManagerInst().AddClass(strName)
	if err != nil {
		return false
	}
	return true
}

func (this *ApiArticleClassController) Post() {
	if !models.TestManagerInst().BgManagerTestEnabled(&this.Controller) {
		if !blogSess.BgManagerEnabled(&this.Controller) {
			this.Abort("404")
			return
		}
	}

	result := false

	cmd := this.GetString("cmd")
	switch cmd {
	case "add":
		result = this.handleAdd()
		break
	}

	if result {
		this.Data["result"] = "success"
	} else {
		this.Data["result"] = "fail"
	}
	this.TplName = "manager/ApiArticleSort.html"
}

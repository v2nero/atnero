package controllers

import (
	"atnero.com/blog/models"
	blogSess "atnero.com/blog/models/session"
	"github.com/astaxie/beego"
)

type CommonPageContainerController struct {
	beego.Controller
}

type articleClassExpose struct {
	Id   int64
	Name string
}

func (this *CommonPageContainerController) InitLayout() {
	this.Layout = "common/PageContainer.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["PageHeader"] = "common/PageHeader.html"
	this.LayoutSections["PageFooter"] = "common/PageFooter.html"
	userInfo := make(map[string]interface{})
	if blogSess.Logined(&this.Controller) {
		usrName, usrId, err := blogSess.GetUserBaseInfo(&this.Controller)
		if err == nil {
			userInfo["name"] = usrName
			userInfo["id"] = usrId
			this.Data["CommonUserInfo"] = userInfo
		}
	}

	if models.TestManagerInst().BgManagerTestEnabled(&this.Controller) ||
		blogSess.BgManagerEnabled(&this.Controller) {
		this.Data["BgManagerEnabled"] = true
	}

	articleClassesList := []articleClassExpose{}
	classList := models.ArticleManagerInst().GetClasses()
	for _, v := range classList {
		classId, err := models.ArticleManagerInst().GetClassId(v)
		if err != nil {
			continue
		}
		c := articleClassExpose{
			Id:   classId,
			Name: v,
		}
		articleClassesList = append(articleClassesList, c)
	}
	this.Data["CommonClassList"] = articleClassesList
}

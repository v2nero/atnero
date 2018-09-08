package controllers

import (
	"github.com/astaxie/beego"
)

type CommonPageContainerController struct {
	beego.Controller
}

func (this *CommonPageContainerController) InitLayout() {
	this.Layout = "common/PageContainer.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["PageHeader"] = "common/PageHeader.html"
	this.LayoutSections["PageFooter"] = "common/PageFooter.html"
}

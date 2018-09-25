package models

import (
	"github.com/astaxie/beego"
	"regexp"
)

type TestManagerIf interface {
	BgManagerTestEnabled(c *beego.Controller) bool
}

type localTestManager struct {
	configBgManagerTestEnabled bool
}

func (this *localTestManager) init() {
	bgMangerTestEnabled, err := beego.AppConfig.Bool("test::bgmanager_test_enable")
	if err != nil {
		this.configBgManagerTestEnabled = false
	} else {
		this.configBgManagerTestEnabled = bgMangerTestEnabled
	}
	if this.configBgManagerTestEnabled {
		beego.Warning("[Test] Background manager test enabled")
	}
}

func (this *localTestManager) isLocalIp(c *beego.Controller) bool {
	request := c.Ctx.Request
	add := request.RemoteAddr
	matched, err := regexp.MatchString(`^::\d{1,3}`, add)
	if err != nil {
		return false
	}
	return matched
}

func (this *localTestManager) BgManagerTestEnabled(
	c *beego.Controller) bool {
	if !this.configBgManagerTestEnabled {
		return false
	}

	return this.isLocalIp(c)
}

var myTestManager TestManagerIf

func init() {
	mng := &localTestManager{}
	mng.init()
	myTestManager = mng
}

func TestManagerInst() TestManagerIf {
	return myTestManager
}

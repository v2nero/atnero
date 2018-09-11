package session

import (
	"atnero.com/blog/models/db"
	"github.com/astaxie/beego"
)

func EnableBgManager(c *beego.Controller, verifyCode string) bool {
	bRet := db.DbMgrInst().VerifyBgManagerPwd(verifyCode)
	if !bRet {
		return bRet
	}

	c.SetSession("BgManagerEnable", true)
	return true
}

func BgManagerEnabled(c *beego.Controller) bool {
	return c.GetSession("BgManagerEnable") != nil
}

func DisableBgManager(c *beego.Controller) {
	c.DelSession("BgManagerEnable")
}

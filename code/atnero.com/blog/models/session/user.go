package session

import (
	"atnero.com/blog/models/db"
	"crypto/md5"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

func Logined(c *beego.Controller) bool {
	return c.GetSession("user") != nil
}

// Login 登陆
// 如果之前已经登陆，则返回失败
// 如果用户名密码验证失败，也返回失败
func Login(c *beego.Controller, user string, pwd string) bool {
	if c.GetSession("user") != nil {
		return false
	}
	var dbUser db.Users
	o := orm.NewOrm()
	err := o.Raw("SELECT * FROM users WHERE name = ?", user).QueryRow(&user)
	if err != nil {
		return false
	}
	strMd5 := fmt.Sprintf("%x", md5.Sum([]byte(pwd)))
	if dbUser.Password != strMd5 {
		return false
	}
	sessUserInfo := make(map[string]interface{})
	sessUserInfo["name"] = dbUser.Name
	sessUserInfo["id"] = dbUser.Id
	sessUserInfo["rightset"] = dbUser.Rightset
	c.SetSession("user", sessUserInfo)
	return true
}

func Logout(c *beego.Controller) {
	c.DelSession("user")
}

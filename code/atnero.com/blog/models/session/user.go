package session

import (
	"atnero.com/blog/models"
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
	err := o.Raw("SELECT * FROM users WHERE name = ?", user).QueryRow(&dbUser)
	if err != nil {
		return false
	}
	strMd5 := fmt.Sprintf("%x", md5.Sum([]byte(pwd)))
	if dbUser.Password != strMd5 {
		return false
	}
	rightSetName, err := models.UserRightsMngInst().GetRightSetNameById(dbUser.Rightset)
	if err != nil {
		return false
	}
	sessUserInfo := make(map[string]interface{})
	sessUserInfo["name"] = dbUser.Name
	sessUserInfo["id"] = dbUser.Id
	sessUserInfo["rightset"] = rightSetName
	c.SetSession("user", sessUserInfo)
	return true
}

func Logout(c *beego.Controller) {
	c.DelSession("user")
}

func UserHasRight(c *beego.Controller, item string) bool {
	var userRightSet string
	for {
		sessUserInfo := c.GetSession("user").(map[string]interface{})
		if sessUserInfo != nil {
			s := sessUserInfo["rightset"]
			if s != nil {
				userRightSet = s.(string)
				if len(userRightSet) != 0 {
					break
				}
			}
		}
		var err error
		userRightSet, err = models.UserRightsMngInst().GetDefaultRightSetName("tourist_rightset")
		if err != nil {
			beego.Error("[RightSet]", err)
			break
		}
		break
	}
	if len(userRightSet) == 0 {
		return false
	}
	return models.UserRightsMngInst().RightSetHasRightItem(userRightSet, item)
}

func init() {
	models.AddDependencyRightSet("tourist_rightset", "user", "网站游客基本权限集")
}

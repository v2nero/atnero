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
		sessUserInfo, ok := c.GetSession("user").(map[string]interface{})
		if ok {
			s := sessUserInfo["rightset"]
			if s != nil {
				userRightSet = s.(string)
				if len(userRightSet) != 0 {
					break
				}
			}
		}
		hasRightSet := models.UserRightsMngInst().HasRightSet("tourist_rightset")
		if !hasRightSet {
			beego.Error("[RightSet] missing 'tourist_rightset'")
			break
		}
		userRightSet = "tourist_rightset"
		break
	}
	if len(userRightSet) == 0 {
		return false
	}
	if !models.UserRightsMngInst().RightSetHasRightItem(userRightSet, item) {
		return false
	}
	enabled, err := models.UserRightsMngInst().RightItemEnabled(item)
	if err != nil {
		return false
	}
	return enabled
}

func GetUserBaseInfo(c *beego.Controller) (name string, id int64, errRet error) {
	if !Logined(c) {
		errRet = fmt.Errorf("not login")
		return
	}
	sessUserInfo := c.GetSession("user").(map[string]interface{})
	name = sessUserInfo["name"].(string)
	id = sessUserInfo["id"].(int64)
	return
}

func init() {
	models.AddDependencyRightSet("tourist_rightset", "user", "网站游客基本权限集")
}

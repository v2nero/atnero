package controllers

import (
	"atnero.com/blog/models"
	"atnero.com/blog/models/db"
	blogSess "atnero.com/blog/models/session"
	"crypto/md5"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

type RegistActionController struct {
	CommonPageContainerController
}

type registInfo struct {
	Name  string `form:"user"`
	Pwd   string `form:"password"`
	Email string `form:"email"`
}

func checkPwd(pwd string) bool {
	bDigit := false
	bLowChar := false
	bUpChar := false
	bSpecial := false
	if len := len(pwd); len < 6 || len > 24 {
		return false
	}
	for _, r := range pwd {
		if r <= 32 || r > 126 {
			return false
		}
		if r >= '0' && r <= '9' {
			bDigit = true
		} else if r >= 'a' && r <= 'z' {
			bLowChar = true
		} else if r >= 'A' && r <= 'z' {
			bUpChar = true
		} else {
			bSpecial = true
		}
	}

	return bDigit && bLowChar && bUpChar && bSpecial
}

func (this *RegistActionController) Post() {
	this.InitLayout()
	this.TplName = "login/regist_action.html"
	this.Data["Title"] = "注册"

	userRightMng := models.UserRightsMngInst()
	if !userRightMng.HasRightSet("base_user_rightset") {
		beego.Error("missing right set ", "base_user_rightset")
		this.Data["InternalError"] = true
		return
	}
	baseRightSetId, err := userRightMng.GetRightSetId("base_user_rightset")
	if err != nil {
		beego.Error(err)
		this.Data["InternalError"] = true
		return
	}

	//已经登陆，显示错误
	if blogSess.Logined(&this.Controller) {
		this.Data["AlreadyLogin"] = true
		return
	}

	this.Data["RegistResult"] = false

	//读取注册信息
	registInfo := registInfo{}
	if err := this.ParseForm(&registInfo); err != nil {
		this.Data["FailReason"] = "解析表单数据出错"
		return
	}

	//用户名不正常
	if len(registInfo.Name) < 4 || len(registInfo.Name) > 24 {
		this.Data["FailReason"] = "用户名长度不在允许范围内"
		return
	}

	//检测密码规则
	if !checkPwd(registInfo.Pwd) {
		this.Data["FailReason"] = "密码不符合规则"
		return
	}

	v := validation.Validation{}
	v.Email(registInfo.Email, "email")
	if v.HasErrors() {
		this.Data["FailReason"] = "无效的邮箱地址"
		return
	}

	strMd5 := fmt.Sprintf("%x", md5.Sum([]byte(registInfo.Pwd)))

	user := db.Users{
		Name:     registInfo.Name,
		Password: strMd5,
		Email:    registInfo.Email,
		Rightset: baseRightSetId,
	}
	o := orm.NewOrm()
	_, err = o.Insert(&user)
	if err != nil {
		beego.Error(err)
		this.Data["FailReason"] = "注册失败"
		return
	}
	this.Data["RegistResult"] = true
}

func init() {
	models.AddDependencyRightSet("base_user_rightset", "Regist Action", "注册用户基本权限集")
}

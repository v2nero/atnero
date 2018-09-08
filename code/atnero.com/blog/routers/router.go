package routers

import (
	"atnero.com/blog/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/manager/userright", &controllers.ManagerUserRightController{})
	beego.Router("/login", &controllers.LoginController{})
}

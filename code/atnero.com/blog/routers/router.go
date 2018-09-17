package routers

import (
	"atnero.com/blog/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/manager/userrightitems", &controllers.ManagerUserRightItemsController{})
	beego.Router("/manager/userrightsets", &controllers.ManagerUserRightSetsController{})
	beego.Router("/manager/userrightset_modify/:setname", &controllers.ManagerUserRightSetModifyController{})
	beego.Router("/manager/defaultuserrightsets", &controllers.ManagerDefaultUserRightSetsController{})
	beego.Router("/manager/api_rightitem", &controllers.ApiRightItemController{})
	beego.Router("/manager/api_rightset", &controllers.ApiRightSetController{})
	beego.Router("/manager/default_rightset", &controllers.ApiDefaultRightSetController{})
	beego.Router("/manager/verify", &controllers.ManagerVerifyController{})
	beego.Router("/login", &controllers.LoginController{})
	beego.Router("/article/create", &controllers.ArticleCreateController{})
}

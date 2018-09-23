package routers

import (
	"atnero.com/blog/controllers"
	"github.com/astaxie/beego"
)

func init() {
	//beego.Router("/", &controllers.MainController{})
	beego.Router("/", &controllers.ArticleListController{})
	beego.Router("/manager/userrightitems", &controllers.ManagerUserRightItemsController{})
	beego.Router("/manager/userrightsets", &controllers.ManagerUserRightSetsController{})
	beego.Router("/manager/userrightset_modify/:setname", &controllers.ManagerUserRightSetModifyController{})
	//beego.Router("/manager/defaultuserrightsets", &controllers.ManagerDefaultUserRightSetsController{})
	beego.Router("/manager/articlesorts", &controllers.ManagerArticleSortsController{})
	beego.Router("/manager/articleclasses", &controllers.ManagerArticleClassesController{})
	beego.Router("/manager/api_rightitem", &controllers.ApiRightItemController{})
	beego.Router("/manager/api_rightset", &controllers.ApiRightSetController{})
	beego.Router("/manager/api_articlesort", &controllers.ApiArticleSortController{})
	beego.Router("/manager/api_articleclass", &controllers.ApiArticleClassController{})
	//beego.Router("/manager/default_rightset", &controllers.ApiDefaultRightSetController{})
	beego.Router("/manager/verify", &controllers.ManagerVerifyController{})
	beego.Router("/login", &controllers.LoginController{})
	beego.Router("/logout", &controllers.LogoutController{})
	beego.Router("/regist", &controllers.RegistController{})
	beego.Router("/regist_action", &controllers.RegistActionController{})
	beego.Router("/article/create", &controllers.ArticleCreateController{})
	beego.Router("/article/view/:id", &controllers.ArticleViewController{})
	beego.Router("/article/list/?:index", &controllers.ArticleListController{})
	beego.Router("/article/listmine/?:index", &controllers.ArticleListMineController{})
	beego.Router("/api_article", &controllers.ApiArticleController{})
}

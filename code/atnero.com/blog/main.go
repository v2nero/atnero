package main

import (
	"atnero.com/blog/db"
	_ "atnero.com/blog/routers"
	"github.com/astaxie/beego"
)

func main() {
	//初始化数据库
	dbConfig := db.DBConfig{
		Host:         "localhost",
		Port:         "3306",
		Database:     "blog",
		Username:     "blog",
		Password:     "123456",
		MaxIdleConns: 100,
		MaxOpenConns: 20,
	}
	db.NewDatabaseManager(&dbConfig)

	//TODO:
	//1. 检查db版本号
	//2. 检查是否开启管理
	//3. 查出权限映射表

	beego.AddTemplateExt("html")
	beego.Run()
}

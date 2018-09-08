package db

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"sync"
)

type DBConfig struct {
	Host         string
	Port         string
	Database     string
	Username     string
	Password     string
	MaxIdleConns int //最大空闲连接
	MaxOpenConns int //最大连接数
}

type DatabaseManager struct {
	DbConfig       *DBConfig
	DbVersion      string
	bgMangerEnable bool
	mutex          sync.Mutex
}

func NewDatabaseManager(dbConfig *DBConfig) *DatabaseManager {
	mng := &DatabaseManager{
		DbConfig: dbConfig,
	}
	mng.init()
	return mng
}

func (mgr *DatabaseManager) init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	ds := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
		mgr.DbConfig.Username,
		mgr.DbConfig.Password,
		mgr.DbConfig.Host,
		mgr.DbConfig.Port,
		mgr.DbConfig.Database)
	logs.Info("datasource=[%s]", ds)
	err := orm.RegisterDataBase("default", "mysql", ds, mgr.DbConfig.MaxIdleConns, mgr.DbConfig.MaxOpenConns)
	if err != nil {
		logs.Error("%+v", err)
		panic(err)
	}
	orm.RegisterModel(new(DbVersion))
	orm.RegisterModel(new(DbManager))
	orm.RegisterModel(new(UserRightItem))
	orm.RegisterModel(new(UserRightSet))
	orm.RegisterModel(new(UserRightSet2itemMap))
	orm.RegisterModel(new(Users))
	orm.RegisterModel(new(ArticleSorts))
	orm.RegisterModel(new(ArticleClasses))
	orm.RegisterModel(new(ArticleLabels))
	orm.RegisterModel(new(Articles))
	orm.RegisterModel(new(ArticleAttachedLabels))
	orm.RegisterModel(new(ArticleComments))
}

func (mgr *DatabaseManager) GetBgManagerEnable() bool {
	return mgr.bgMangerEnable
}

func (mgr *DatabaseManager) SetBgManagerEnable(enable bool) {
	var err error
	mgr.mutex.Lock()

	var dbMng DbManager
	o := orm.NewOrm()
	qs := o.QueryTable("dbmanager")
	err = qs.One(&dbMng)
	if err != nil {
		mgr.mutex.Unlock()
		panic("no records in dbmanager")
	}
	dbMng.Enable = enable
	num, err := o.Update(&dbMng)
	if err != nil || num < 1 {
		if err != nil {
			logs.Error("update dbmanager with error: %+v", err)
		} else {
			logs.Error("update dbmanager fail, update number is %d", num)
		}
		mgr.mutex.Unlock()
		panic("fail to update dbmanger table")
	}
	mgr.bgMangerEnable = enable
	mgr.mutex.Unlock()
}

func (mgr *DatabaseManager) loadRights() {
	//TODO:
}

var dbMgr *DatabaseManager

func init() {
	var dbConfig DBConfig
	dbConfig.Host = beego.AppConfig.String("mysqlhost")
	dbConfig.Port = beego.AppConfig.String("mysqlport")
	dbConfig.Database = beego.AppConfig.String("mysqldb")
	dbConfig.Username = beego.AppConfig.String("mysqluser")
	dbConfig.Password = beego.AppConfig.String("mysqlpwd")
	var err error
	dbConfig.MaxIdleConns, err = beego.AppConfig.Int("mysqlmaxidleconns")
	if err != nil {
		dbConfig.MaxIdleConns = 100
	}
	dbConfig.MaxOpenConns, err = beego.AppConfig.Int("mysqlmaxopenconns")
	if err != nil {
		dbConfig.MaxOpenConns = 10
	}
	dbMgr = NewDatabaseManager(&dbConfig)
}

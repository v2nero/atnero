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
	DbConfig    *DBConfig
	DbVersion   string
	bgMngEnable bool
	mutex       sync.Mutex
}

func (mng *DatabaseManager) GetBgManagerEnable() bool {
	return mng.bgMngEnable
}

func (mng *DatabaseManager) SetBgManagerEnable(enable bool) error {
	var err error
	mng.mutex.Lock()

	for {
		var bgMng BgManagerEnable
		o := orm.NewOrm()
		qs := o.QueryTable("dbmanager")
		err = qs.One(&bgMng)
		if err != nil {
			logs.Error("query table dbmanager return with error: %+v", err)
			return err
		}
		bgMng.Enable = enable
		num, err := o.Update(&bgMng)
		if err != nil || num < 1 {
			if err != nil {
				logs.Error("update dbmanager with error: %+v", err)
			} else {
				logs.Error("update dbmanager fail, update number is %d", num)
			}

			err = fmt.Errorf("fail to update dbmanger table")
			break
		}
		mng.bgMngEnable = enable
		break
	}

	mng.mutex.Unlock()
	return err
}

func (mng *DatabaseManager) initDbVersion() {
	var version Version
	o := orm.NewOrm()
	if err := o.Read(&version); err != nil {
		panic(err)
	}
	mng.DbVersion = version.Version
}

func (mng *DatabaseManager) initDbMngEnable() {
	var en BgManagerEnable
	o := orm.NewOrm()
	if err := o.Read(&en); err != nil {
		panic(err)
	}
	mng.bgMngEnable = en.Enable
}

func (mng *DatabaseManager) init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	ds := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
		mng.DbConfig.Username,
		mng.DbConfig.Password,
		mng.DbConfig.Host,
		mng.DbConfig.Port,
		mng.DbConfig.Database)
	logs.Info("datasource=[%s]", ds)
	err := orm.RegisterDataBase("default", "mysql", ds, mng.DbConfig.MaxIdleConns, mng.DbConfig.MaxOpenConns)
	if err != nil {
		logs.Error("%+v", err)
		panic(err)
	}
	orm.RegisterModel(new(Version))
	orm.RegisterModel(new(BgManagerEnable))
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

	mng.initDbVersion()
	mng.initDbMngEnable()
}

func newDatabaseManager(dbConfig *DBConfig) *DatabaseManager {
	mng := &DatabaseManager{
		DbConfig: dbConfig,
	}
	mng.init()
	return mng
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
		dbConfig.MaxIdleConns = 10
	}
	dbConfig.MaxOpenConns, err = beego.AppConfig.Int("mysqlmaxopenconns")
	if err != nil {
		dbConfig.MaxOpenConns = 10
	}
	dbMgr = newDatabaseManager(&dbConfig)
}

func DbMgrInst() *DatabaseManager {
	return dbMgr
}

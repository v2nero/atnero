package db

import (
	"time"
)

//文章分类，如原创，转发等
type ArticleSorts struct {
	Id          int64 `orm:"auto"`
	Name        string
	Discription string `orm:"column(dsc)"`
}

//文章用户分类
type ArticleClasses struct {
	Id          int64 `orm:"auto"`
	UserId      int64
	Name        string
	Discription string `orm:"column(dsc)"`
}

type ArticleLabels struct {
	Id          int64 `orm:"auto"`
	UserId      int64
	Name        string
	Discription string `orm:"column(dsc)"`
}

type Articles struct {
	Id             int64 `orm:"auto"`
	UserId         int64
	Title          string
	ShortDsc       string
	Content        string
	SortId         int64
	ClassId        int64
	Published      bool
	CreateTime     time.Time `orm:"auto_now_add;type(datetime)"`
	LastupdateTime time.Time `orm:"auto_now;type(datatime)"`
	ViewCount      int64
}

type ArticleAttachedLabels struct {
	Id        int64 `orm:"auto"`
	ArticleId int64
	LabelId   int64
}

type ArticleComments struct {
	Id        int64 `orm:"auto"`
	ArticleId int64
	Content   string
}

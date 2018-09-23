package dbview

import (
	"time"
)

type ArticleCommentView struct {
	Id         int64
	ArticleId  int64
	UserName   string
	Email      string
	Content    string
	CreateTime time.Time
}

type ArticleLabelView struct {
	Id   int64 //attached label record id
	Name string
}

type ArticleDataView struct {
	Id             int64
	UserId         int64
	UserName       string
	Title          string
	Content        string
	SortName       string
	ClassName      string
	Published      bool
	CreateTime     time.Time
	LastupdateTime time.Time
	Labels         *[]ArticleLabelView
	Comments       *[]ArticleCommentView
}

type AritcleShortView struct {
	Id        int64
	Title     string
	ViewCount int64
}

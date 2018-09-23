package models

import (
	"atnero.com/blog/models/db"
	"atnero.com/blog/models/dbview"
)

type ArticleManagerIf interface {
	Init() error
	GetSorts() []string
	GetSortId(name string) (int64, error)
	AddSort(name string) error
	GetSortName(id int64) (string, error)
	//No update or delete interface

	GetClasses() []string
	GetClassId(name string) (int64, error)
	AddClass(name string) error
	GetClassName(id int64) (string, error)
	//no update or delete interface

	GetLabelsByUserId(id int64) []string
	GetLabelsByArticleId(id int64) []string
	AddLabel(userId int64, name string) error
	AttachLabelToArticle(name string, id int64) error

	AddComment(articleId int64, userName string, email string, content string) (int64, error)
	DeleteComment(id int64) error
	//no udpate or delete interface

	AddArticle(userId int64, title string, content string,
		sortName string, className string, publish bool) (int64, error)
	UpdateArticle(id int64, title string, content string,
		sortName string, className string, publish bool) error
	SetArticlePublishState(id int64, publish bool) error
	GetArticlePublisState(id int64) (bool, error)
	GetArticlesNumOfUser(userId int64, publishedOnly bool) (int, error)
	GetArticlesShortViewOfUser(userId int64, publishedOnly bool,
		index int64, limit int64) ([]dbview.AritcleShortView, error)
	GetArticlesNumOfAll(publishedOnly bool) (int, error)
	GetArticlesShortViewOfAll(publishedOnly bool,
		index int64, limit int64) ([]dbview.AritcleShortView, error)
	GetArticleShortView(id int64) (*dbview.AritcleShortView, error)
	GetArticleData(id int64) (*dbview.ArticleDataView, error)
}

var myArticleManagerInst ArticleManagerIf

func init() {
	myArticleManagerInst = new(db.DbArticleManager)
	myArticleManagerInst.Init()
}

func ArticleManagerInst() ArticleManagerIf {
	return myArticleManagerInst
}

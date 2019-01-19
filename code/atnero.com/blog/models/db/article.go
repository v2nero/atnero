package db

import (
	"atnero.com/blog/models/dbview"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"sort"
	"sync"
	"time"
)

//文章分类，如原创，转发等
type ArticleSorts struct {
	Id   int64 `orm:"auto"`
	Name string
}

//文章用户分类
type ArticleClasses struct {
	Id   int64 `orm:"auto"`
	Name string
}

type ArticleLabels struct {
	Id     int64 `orm:"auto"`
	UserId int64
	Name   string
}

type Articles struct {
	Id             int64 `orm:"auto"`
	UserId         int64
	Title          string
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
	Id         int64 `orm:"auto"`
	ArticleId  int64
	UserName   string
	Email      string
	Content    string
	CreateTime time.Time `orm:"auto_now;type(datetime)"`
}

//==========================MANAGER=========================

type DbArticleManager struct {
	sortMap    map[string]*ArticleSorts
	sortIdMap  map[int64]*ArticleSorts
	classMap   map[string]*ArticleClasses
	classIdMap map[int64]*ArticleClasses
	mutex      sync.Mutex
	//labels, comments, article暂时不进行缓存
}

//仅供init使用
func (this *DbArticleManager) loadSorts() {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	var sorts []*ArticleSorts
	o := orm.NewOrm()
	qs := o.QueryTable("article_sorts")
	num, err := qs.All(&sorts)
	if num == 0 {
		return
	}
	if err != nil {
		panic(err)
	}
	for _, i := range sorts {
		this.sortMap[i.Name] = i
		this.sortIdMap[i.Id] = i
	}
}

//仅供init使用
func (this *DbArticleManager) loadClasses() {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	var classes []*ArticleClasses
	o := orm.NewOrm()
	qs := o.QueryTable("article_classes")
	num, err := qs.All(&classes)
	if num == 0 {
		return
	}
	if err != nil {
		panic(err)
	}
	for _, i := range classes {
		this.classMap[i.Name] = i
		this.classIdMap[i.Id] = i
	}
}

func (this *DbArticleManager) Init() error {
	this.sortMap = make(map[string]*ArticleSorts)
	this.sortIdMap = make(map[int64]*ArticleSorts)
	this.classMap = make(map[string]*ArticleClasses)
	this.classIdMap = make(map[int64]*ArticleClasses)
	this.loadSorts()
	this.loadClasses()
	return nil
}

func (this *DbArticleManager) GetSorts() []string {
	var sorts []string
	this.mutex.Lock()
	for k := range this.sortMap {
		sorts = append(sorts, k)
	}
	this.mutex.Unlock()
	sort.Strings(sorts)
	return sorts
}

func (this *DbArticleManager) GetSortId(name string) (int64, error) {
	var idRet int64
	var errRet error
	this.mutex.Lock()
	for {
		n, b := this.sortMap[name]
		if !b {
			errRet = fmt.Errorf("no such sort %s", name)
			break
		}
		idRet = n.Id
		break
	}
	this.mutex.Unlock()
	return idRet, errRet
}

func (this *DbArticleManager) GetSortName(id int64) (string, error) {
	var name string
	var errRet error
	this.mutex.Lock()
	for {
		n, b := this.sortIdMap[id]
		if !b {
			errRet = fmt.Errorf("no such sort %v", id)
			break
		}
		name = n.Name
		break
	}
	this.mutex.Unlock()
	return name, errRet
}

//非线程安全
func (this *DbArticleManager) ormAddSort(
	name string) (int64, error) {
	var sort ArticleSorts
	sort.Name = name

	o := orm.NewOrm()
	id, err := o.Insert(&sort)
	if err != nil {
		logs.Error("[ORM]insert article sort item fail with error:", err)
	}
	return id, err
}

func (this *DbArticleManager) AddSort(name string) error {
	var errRet error
	this.mutex.Lock()
	for {
		_, b := this.sortMap[name]
		if b {
			errRet = fmt.Errorf("Sort [%s] already exist", name)
			break
		}
		id, err := this.ormAddSort(name)
		if err != nil {
			break
		}
		sort := &ArticleSorts{
			Id:   id,
			Name: name,
		}
		this.sortMap[name] = sort
		this.sortIdMap[id] = sort
		break
	}
	this.mutex.Unlock()
	return errRet
}

func (this *DbArticleManager) GetClasses() []string {
	var class []string
	this.mutex.Lock()
	for k := range this.classMap {
		class = append(class, k)
	}
	this.mutex.Unlock()
	sort.Strings(class)
	return class
}

func (this *DbArticleManager) GetClassId(name string) (int64, error) {
	var idRet int64
	var errRet error
	this.mutex.Lock()
	for {
		n, b := this.classMap[name]
		if !b {
			errRet = fmt.Errorf("no such class %s", name)
			break
		}
		idRet = n.Id
		break
	}
	this.mutex.Unlock()
	return idRet, errRet
}

func (this *DbArticleManager) GetClassName(id int64) (string, error) {
	var name string
	var errRet error
	this.mutex.Lock()
	for {
		n, b := this.classIdMap[id]
		if !b {
			errRet = fmt.Errorf("no such class %v", id)
			break
		}
		name = n.Name
		break
	}
	this.mutex.Unlock()
	return name, errRet
}

//非线程安全
func (this *DbArticleManager) ormAddClass(
	name string) (int64, error) {
	var class ArticleClasses
	class.Name = name

	o := orm.NewOrm()
	id, err := o.Insert(&class)
	if err != nil {
		logs.Error("[ORM]insert article class item fail with error:%+v", err)
	}
	return id, err
}

func (this *DbArticleManager) AddClass(name string) error {
	var errRet error
	this.mutex.Lock()
	for {
		_, b := this.classMap[name]
		if b {
			errRet = fmt.Errorf("Class [%s] already exist", name)
			break
		}
		id, err := this.ormAddClass(name)
		if err != nil {
			break
		}
		c := &ArticleClasses{
			Id:   id,
			Name: name,
		}
		this.classMap[name] = c
		this.classIdMap[id] = c
		break
	}
	this.mutex.Unlock()
	return errRet
}

func (this *DbArticleManager) GetLabelsByUserId(id int64) []string {
	var labels []string
	//this.mutex.Lock()
	for {
		tlabels := []struct {
			Name string
		}{}
		o := orm.NewOrm()
		_, err := o.Raw("select name from article_labels where user_id = ?", id).QueryRows(&tlabels)
		if err != nil {
			break
		}
		for _, v := range tlabels {
			labels = append(labels, v.Name)
		}
		break
	}
	//this.mutex.Unlock()
	return labels
}

func (this *DbArticleManager) GetLabelsByArticleId(id int64) []string {
	var labels []string
	//this.mutex.Lock()
	for {
		tlabels := []struct {
			Name string
		}{}
		o := orm.NewOrm()
		_, err := o.Raw(`SELECT l.name FROM article_attached_labels a INNER JOIN \
				article_labels l ON a.label_id = l.id where a.article_id = ?`, id).QueryRows(&tlabels)
		if err != nil {
			break
		}
		for _, v := range tlabels {
			labels = append(labels, v.Name)
		}
		break
	}
	//this.mutex.Unlock()
	return labels
}

func (this *DbArticleManager) AddLabel(userId int64, name string) error {
	var errRet error
	for {
		label := ArticleLabels{
			UserId: userId,
			Name:   name,
		}
		o := orm.NewOrm()
		_, errRet = o.Insert(&label)
		if errRet != nil {
			break
		}
		break
	}
	return errRet
}

func (this *DbArticleManager) AttachLabelToArticle(name string, id int64) error {
	var errRet error
	for {
		o := orm.NewOrm()
		res, err := o.Raw(`INSERT INTO article_attached_labels(article_id, label_id) \
				SELECT ?, l.name FROM article_labels l INNER JOIN \
				articles a ON l.user_id = a.user_id WHERE a.id = ?`, id, id).Exec()
		if err != nil {
			errRet = err
			break
		}
		if num, err := res.RowsAffected(); num == 0 || err != nil {
			errRet = fmt.Errorf("no rows inserted")
			break
		}
		break
	}
	return errRet
}

func (this *DbArticleManager) AddComment(
	articleId int64, userName string,
	email string, content string) (int64, error) {
	var errRet error
	var idRet int64
	for {
		if len(content) == 0 ||
			len(userName) == 0 ||
			len(email) == 0 {
			errRet = fmt.Errorf("invalid parameter")
			break
		}
		v := validation.Validation{}
		v.Email(&email, "email")
		if v.HasErrors() {
			errRet = fmt.Errorf("fail to validate email address")
			break
		}
		comment := ArticleComments{
			ArticleId: articleId,
			UserName:  userName,
			Email:     email,
			Content:   content,
		}
		o := orm.NewOrm()
		idRet, errRet = o.Insert(&comment)
		if errRet != nil {
			break
		}
		break
	}
	return idRet, errRet
}

func (this *DbArticleManager) DeleteComment(id int64) error {
	comment := ArticleComments{
		Id: id,
	}
	o := orm.NewOrm()
	num, err := o.Delete(&comment)
	if err != nil {
		return err
	}
	if num == 0 {
		return fmt.Errorf("no rows impact")
	}
	return nil
}

func (this *DbArticleManager) AddArticle(userId int64, title string, content string,
	sortName string, className string, publish bool) (int64, error) {
	var idRet int64
	var errRet error
	for {
		sortId, err := this.GetSortId(sortName)
		if err != nil {
			errRet = err
			break
		}
		classId, err := this.GetClassId(className)
		if err != nil {
			errRet = err
			break
		}
		a := Articles{
			UserId:    userId,
			Title:     title,
			Content:   content,
			SortId:    sortId,
			ClassId:   classId,
			Published: publish,
			ViewCount: 0,
		}

		o := orm.NewOrm()
		idRet, errRet = o.Insert(&a)
		if errRet != nil {
			break
		}

		break
	}

	return idRet, errRet
}
func (this *DbArticleManager) UpdateArticle(id int64, title string, content string,
	sortName string, className string, publish bool) error {
	var errRet error
	for {
		sortId, err := this.GetSortId(sortName)
		if err != nil {
			errRet = err
			break
		}
		classId, err := this.GetClassId(className)
		if err != nil {
			errRet = err
			break
		}
		a := Articles{
			Id: id,
		}

		o := orm.NewOrm()
		errRet = o.Read(&a)
		if errRet != nil {
			break
		}

		a.Title = title
		a.Content = content
		a.SortId = sortId
		a.ClassId = classId
		a.Published = publish
		num, err := o.Update(&a)
		if err != nil {
			errRet = err
			break
		}
		if num != 1 {
			errRet = fmt.Errorf("no article updated")
			break
		}
		break
	}

	return errRet
}

func (this *DbArticleManager) SetArticlePublishState(id int64, publish bool) error {
	var errRet error
	for {
		a := Articles{
			Id: id,
		}

		o := orm.NewOrm()
		errRet = o.Read(&a)
		if errRet != nil {
			break
		}

		a.Published = publish
		num, err := o.Update(&a)
		if err != nil {
			errRet = err
			break
		}
		if num != 1 {
			errRet = fmt.Errorf("no article updated")
			break
		}
		break
	}

	return errRet
}

func (this *DbArticleManager) GetArticlePublisState(id int64) (bool, error) {
	var published bool
	var errRet error
	for {
		a := Articles{
			Id: id,
		}

		o := orm.NewOrm()
		errRet = o.Read(&a)
		if errRet != nil {
			break
		}
		published = a.Published
		break
	}

	return published, errRet
}

func (this *DbArticleManager) GetArticlesNumOfUser(
	userId int64, publishedOnly bool) (int, error) {
	var recordCnt int
	var errRet error
	for {
		cntView := struct {
			Cnt int
		}{}

		o := orm.NewOrm()
		if publishedOnly {
			err := o.Raw(`select COUNT(id) as "cnt" from articles where user_id=? AND published=?`,
				userId, true).QueryRow(&cntView)
			if err != nil {
				errRet = err
			}
		} else {
			err := o.Raw(`select COUNT(id) as "cnt" from articles where user_id=?`,
				userId).QueryRow(&cntView)
			if err != nil {
				errRet = err
			}
		}
		recordCnt = cntView.Cnt
		break
	}

	return recordCnt, errRet
}

func (this *DbArticleManager) GetArticleShortView(id int64) (*dbview.AritcleShortView, error) {
	var view dbview.AritcleShortView
	var errRet error
	for {
		o := orm.NewOrm()

		err := o.Raw("select a.id, a.user_id, users.name as 'user_name', a.title, sort.name as 'sort_name', class.name as 'class_name', a.published, a.view_count, a.lastupdate_time from articles a INNER JOIN users ON a.user_id = users.id INNER JOIN article_sorts sort ON sort.id = a.sort_id INNER JOIN article_classes class ON class.id = a.class_id where a.id=?",
			id).QueryRow(&view)
		if err != nil {
			errRet = err
		}

		break
	}
	return &view, errRet
}

func (this *DbArticleManager) GetArticlesShortViewOfUser(
	userId int64, publishedOnly bool,
	index int64, limit int64) ([]dbview.AritcleShortView, error) {
	var views []dbview.AritcleShortView
	var errRet error
	for {
		o := orm.NewOrm()
		if publishedOnly {
			_, err := o.Raw("select a.id, a.user_id, users.name as 'user_name', a.title, sort.name as 'sort_name', class.name as 'class_name', a.published, a.view_count, a.lastupdate_time from articles a INNER JOIN users ON a.user_id = users.id INNER JOIN article_sorts sort ON sort.id = a.sort_id INNER JOIN article_classes class ON class.id = a.class_id where a.user_id=? AND a.published=? limit ?,?",
				userId, true, index*limit, limit).QueryRows(&views)
			if err != nil {
				errRet = err
			}
		} else {
			_, err := o.Raw("select a.id, a.user_id, users.name as 'user_name', a.title, sort.name as 'sort_name', class.name as 'class_name', a.published, a.view_count, a.lastupdate_time from articles a INNER JOIN users ON a.user_id = users.id INNER JOIN article_sorts sort ON sort.id = a.sort_id INNER JOIN article_classes class ON class.id = a.class_id where a.user_id=? limit ?,?",
				userId, index*limit, limit).QueryRows(&views)
			if err != nil {
				errRet = err
			}
		}

		break
	}
	return views, errRet
}

func (this *DbArticleManager) GetArticlesNumOfClass(classId int64) (int, error) {
	var recordCnt int
	var errRet error
	for {
		cntView := struct {
			Cnt int
		}{}

		o := orm.NewOrm()
		err := o.Raw(`select COUNT(id) as "cnt" from articles where class_id=? AND published=true`,
			classId).QueryRow(&cntView)
		if err != nil {
			errRet = err
		}
		recordCnt = cntView.Cnt
		break
	}

	return recordCnt, errRet
}

func (this *DbArticleManager) GetArticlesShortViewOfClass(classId int64,
	index int64, limit int64) ([]dbview.AritcleShortView, error) {
	var views []dbview.AritcleShortView
	var errRet error
	for {
		o := orm.NewOrm()
		_, err := o.Raw("select a.id, a.user_id, users.name as 'user_name', a.title, sort.name as 'sort_name', class.name as 'class_name', a.published, a.view_count, a.lastupdate_time from articles a INNER JOIN users ON a.user_id = users.id INNER JOIN article_sorts sort ON sort.id = a.sort_id INNER JOIN article_classes class ON class.id = a.class_id where a.class_id=? AND a.published=true  order by a.id desc limit ?,?",
			classId, index*limit, limit).QueryRows(&views)
		if err != nil {
			errRet = err
		}

		break
	}
	return views, errRet
}

func (this *DbArticleManager) GetArticlesNumOfAll(
	publishedOnly bool) (int, error) {
	var recordCnt int
	var errRet error
	for {
		cntView := struct {
			Cnt int
		}{}

		o := orm.NewOrm()
		if publishedOnly {
			err := o.Raw(`select COUNT(id) as "cnt" from articles where published=?`, true).QueryRow(&cntView)
			if err != nil {
				errRet = err
			}
		} else {
			err := o.Raw(`select COUNT(id) as "cnt" from articles`).QueryRow(&cntView)
			if err != nil {
				errRet = err
			}
		}
		recordCnt = cntView.Cnt
		break
	}

	return recordCnt, errRet
}

func (this *DbArticleManager) GetArticlesShortViewOfAll(publishedOnly bool,
	index int64, limit int64) ([]dbview.AritcleShortView, error) {
	var views []dbview.AritcleShortView
	var errRet error
	for {
		o := orm.NewOrm()
		if publishedOnly {
			_, err := o.Raw("select a.id, a.user_id, users.name as 'user_name', a.title, sort.name as 'sort_name', class.name as 'class_name', a.published, a.view_count, a.lastupdate_time from articles a INNER JOIN users ON a.user_id = users.id INNER JOIN article_sorts sort ON sort.id = a.sort_id INNER JOIN article_classes class ON class.id = a.class_id where published=? order by a.id desc limit ?,?",
				true, index*limit, limit).QueryRows(&views)
			if err != nil {
				errRet = err
			}
		} else {
			_, err := o.Raw("select a.id, a.user_id, users.name as 'user_name', a.title, sort.name as 'sort_name', class.name as 'class_name', a.published, a.view_count, a.lastupdate_time from articles a INNER JOIN users ON a.user_id = users.id INNER JOIN article_sorts sort ON sort.id = a.sort_id INNER JOIN article_classes class ON class.id = a.class_id  order by a.id desc limit ?,?",
				index*limit, limit).QueryRows(&views)
			if err != nil {
				errRet = err
			}
		}
		break
	}
	return views, errRet
}

func (this *DbArticleManager) GetArticleData(id int64) (*dbview.ArticleDataView, error) {
	var view dbview.ArticleDataView
	var errRet error
	for {
		a := Articles{}
		o := orm.NewOrm()
		err := o.Raw("SELECT * FROM articles WHERE id = ?", id).QueryRow(&a)
		if err != nil {
			errRet = err
			break
		}

		sortName, err := this.GetSortName(a.SortId)
		if err != nil {
			errRet = err
			break
		}

		className, err := this.GetClassName(a.ClassId)
		if err != nil {
			errRet = err
			break
		}

		usr := Users{}
		err = o.Raw("SELECT * FROM users WHERE id = ?", a.UserId).QueryRow(&usr)
		if err != nil {
			errRet = err
			break
		}

		var labels []dbview.ArticleLabelView
		_, err = o.Raw("SELECT l.id, l.name FROM article_attached_labels a LEFT OUTER JOIN article_labels l ON a.label_id = l.id  WHERE a.article_id = ?",
			a.Id).QueryRows(&labels)
		if err != nil {
			errRet = err
			break
		}

		var comments []ArticleComments
		_, err = o.Raw("SELECT * FROM article_comments WHERE article_id = ?", a.Id).QueryRows(&comments)
		if err != nil {
			errRet = err
			break
		}

		view.Id = a.Id
		view.ClassName = className
		view.SortName = sortName
		view.Content = a.Content
		view.Title = a.Title
		view.UserId = a.UserId
		view.UserName = usr.Name
		view.Published = a.Published
		view.CreateTime = a.CreateTime
		view.LastupdateTime = a.LastupdateTime
		view.ViewCount = a.ViewCount

		labelViews := []dbview.ArticleLabelView{}
		for _, v := range labels {
			lv := dbview.ArticleLabelView{
				Id:   v.Id,
				Name: v.Name,
			}
			labelViews = append(labelViews, lv)
		}
		view.Labels = &labelViews

		commentViews := []dbview.ArticleCommentView{}
		for _, v := range comments {
			cv := dbview.ArticleCommentView{
				Id:         v.Id,
				ArticleId:  v.ArticleId,
				UserName:   v.UserName,
				Email:      v.Email,
				Content:    v.Content,
				CreateTime: v.CreateTime,
			}
			commentViews = append(commentViews, cv)
		}
		view.Comments = &commentViews
		break
	}
	return &view, errRet
}

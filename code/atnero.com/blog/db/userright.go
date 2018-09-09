package db

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"sync"
)

// ============ORM table mapping===========
type UserRightItem struct {
	Id          int64 `orm:"auto"`
	Name        string
	Enable      bool
	Discription string `orm:"column(dsc)"`
}

type UserRightSet struct {
	Id          int64 `orm:"auto"`
	Name        string
	Discription string `orm:"column(dsc)"`
}

type UserRightSet2itemMap struct {
	Id     int64 `orm:"auto"`
	SetId  int64
	ItemId int64
}

// ============User rights manager instance===========
type rightSet2ItemPair struct {
	setId  int64
	itemId int64
}

type DbUserRightsManager struct {
	itemMap map[string]UserRightItem
	setMap  map[string]UserRightSet
	linkMap map[rightSet2ItemPair]UserRightSet2itemMap
	mutex   sync.Mutex
}

//Init 初始化
//装数据库数据
func (this *DbUserRightsManager) Init() {
	this.itemMap = make(map[string]UserRightItem)
	this.setMap = make(map[string]UserRightSet)
	this.linkMap = make(map[rightSet2ItemPair]UserRightSet2itemMap)
	this.loadRightItems()
	this.loadRightSets()
	this.loadRightSet2ItemMapping()
}

//仅供init使用
func (this *DbUserRightsManager) loadRightItems() {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	var rightItems []*UserRightItem
	o := orm.NewOrm()
	qs := o.QueryTable("user_right_item")
	_, err := qs.All(&rightItems)
	if err != nil {
		panic(err)
	}
	for _, i := range rightItems {
		this.itemMap[i.Name] = *i
	}
}

//仅供init使用
func (this *DbUserRightsManager) loadRightSets() {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	var rightSets []*UserRightSet
	o := orm.NewOrm()
	qs := o.QueryTable("user_right_set")
	_, err := qs.All(&rightSets)
	if err != nil {
		panic(err)
	}
	for _, i := range rightSets {
		this.setMap[i.Name] = *i
	}
}

func (this *DbUserRightsManager) checkRightItemExistanceWithIdWithoutLock(
	itemId int64) bool {
	for _, v := range this.itemMap {
		if v.Id == itemId {
			return true
		}
	}
	return false
}

func (this *DbUserRightsManager) checkRightSetExistanceWithIdWithoutLock(
	setId int64) bool {
	for _, v := range this.setMap {
		if v.Id == setId {
			return true
		}
	}
	return false
}

//仅供init使用
func (this *DbUserRightsManager) loadRightSet2ItemMapping() {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	var mapNodes []*UserRightSet2itemMap
	o := orm.NewOrm()
	qs := o.QueryTable("user_right_set2item_map")
	_, err := qs.All(&mapNodes)
	if err != nil {
		panic(err)
	}
	for _, i := range mapNodes {
		itemExist := this.checkRightItemExistanceWithIdWithoutLock(i.ItemId)
		setExist := this.checkRightSetExistanceWithIdWithoutLock(i.SetId)
		if !itemExist || !setExist {
			panic(fmt.Sprintf(
				"RightItem %d (%v?do:not) exist or RightSet %d (%v?do:not) exist",
				i.ItemId, itemExist,
				i.SetId, setExist))
		}
		link := rightSet2ItemPair{
			setId:  i.SetId,
			itemId: i.ItemId,
		}
		this.linkMap[link] = *i
	}
}

func (this *DbUserRightsManager) GetRightItems() []string {
	var keys []string
	this.mutex.Lock()
	for k, _ := range this.itemMap {
		keys = append(keys, k)
	}
	this.mutex.Unlock()
	return keys
}

func (this *DbUserRightsManager) GetRightSets() []string {
	var keys []string
	this.mutex.Lock()
	for k, _ := range this.setMap {
		keys = append(keys, k)
	}
	this.mutex.Unlock()
	return keys
}

func (this *DbUserRightsManager) HasRightItem(item string) bool {
	var b bool
	this.mutex.Lock()
	_, b = this.itemMap[item]
	this.mutex.Unlock()
	return b
}

func (this *DbUserRightsManager) GetRightItemDiscription(
	item string) (string, bool) {
	var b bool
	this.mutex.Lock()
	r, b := this.itemMap[item]
	this.mutex.Unlock()
	return r.Discription, b
}

func (this *DbUserRightsManager) HasRightSet(set string) bool {
	var b bool
	this.mutex.Lock()
	_, b = this.setMap[set]
	this.mutex.Unlock()
	return b
}

func (this *DbUserRightsManager) GetRightSetDiscription(
	set string) (string, bool) {
	var b bool
	this.mutex.Lock()
	r, b := this.setMap[set]
	this.mutex.Unlock()
	return r.Discription, b
}

func (this *DbUserRightsManager) RightSetHasRightItem(
	set string, item string) bool {
	var bRet bool
	var setExist, itemExist bool
	var setItem UserRightSet
	var itemItem UserRightItem

	bRet = true
	this.mutex.Lock()
	for {

		setItem, setExist = this.setMap[set]
		itemItem, itemExist = this.itemMap[item]
		if !setExist || !itemExist {
			bRet = false
			break
		}
		linkItem := rightSet2ItemPair{
			setId:  setItem.Id,
			itemId: itemItem.Id,
		}
		_, linkExist := this.linkMap[linkItem]
		if !linkExist {
			bRet = false
		}

		break
	}
	this.mutex.Unlock()
	return bRet
}

//非线程安全
func (this *DbUserRightsManager) ormAddRightItem(
	name string, dsc string) (int64, error) {
	var item UserRightItem
	item.Name = name
	item.Discription = dsc
	item.Enable = false

	o := orm.NewOrm()
	id, err := o.Insert(&item)
	if err != nil {
		logs.Error("[ORM]insert user right item fail with error:%+v", err)
	}
	return id, err
}

func (this *DbUserRightsManager) AddRightItem(
	item string,
	dsc string) error {
	var errRet error
	this.mutex.Lock()
	for {
		_, b := this.itemMap[item]
		if b {
			errRet = fmt.Errorf("Right item [%s] already exists", item)
			break
		}
		id, err := this.ormAddRightItem(item, dsc)
		if err != nil {
			break
		}
		rightItem := UserRightItem{
			id, item, false, dsc,
		}
		this.itemMap[item] = rightItem
		break
	}
	this.mutex.Unlock()
	return errRet
}

func (this *DbUserRightsManager) ormUpdateRightItemRecord(
	r *UserRightItem) error {
	o := orm.NewOrm()
	num, err := o.Update(r)
	if num != 1 {
		return fmt.Errorf("update UserRightItem but no rows impacted")
	}
	return err
}

func (this *DbUserRightsManager) EnableRightItem(
	item string,
	enable bool) error {
	var errRet error
	this.mutex.Lock()
	for {
		r, exist := this.itemMap[item]
		if !exist {
			errRet = fmt.Errorf("item %v not exist", item)
			break
		}
		r.Enable = enable
		errRet = this.ormUpdateRightItemRecord(&r)
		break
	}
	this.mutex.Unlock()
	return errRet
}

//非线程安全
func (this *DbUserRightsManager) ormAddRightSet(
	name string, dsc string) (int64, error) {
	var set UserRightSet
	set.Name = name
	set.Discription = dsc

	o := orm.NewOrm()
	id, err := o.Insert(&set)
	if err != nil {
		logs.Error("[ORM]insert user right set fail with error:%+v", err)
	}
	return id, err
}

func (this *DbUserRightsManager) AddRightSet(
	set string,
	dsc string) error {
	var errRet error
	this.mutex.Lock()
	for {
		_, b := this.setMap[set]
		if b {
			errRet = fmt.Errorf("Right set [%s] already exists", set)
			break
		}
		id, err := this.ormAddRightSet(set, dsc)
		if err != nil {
			break
		}
		setItem := UserRightSet{
			id, set, dsc,
		}
		this.setMap[set] = setItem
		break
	}
	this.mutex.Unlock()
	return errRet
}

//非线程安全
func (this *DbUserRightsManager) ormAddRightItem2RightSet(
	itemId int64, setId int64) (int64, error) {

	mapNode := UserRightSet2itemMap{
		ItemId: itemId,
		SetId:  setId,
	}

	o := orm.NewOrm()
	id, err := o.Insert(&mapNode)
	if err != nil {
		logs.Error("[ORM]insert Set2ItemMap node fail with error:%+v", err)
	}
	return id, err
}

func (this *DbUserRightsManager) AddRightItem2RightSet(
	set string, item string) error {
	var errRet error
	this.mutex.Lock()
	for {
		setItem, setExist := this.setMap[set]
		itemItem, itemExist := this.itemMap[item]
		if !setExist || !itemExist {
			errRet = fmt.Errorf("set or item not exist(%v, %v)", setExist, itemExist)
			break
		}
		linkItem := rightSet2ItemPair{
			setId:  setItem.Id,
			itemId: itemItem.Id,
		}
		_, linkExist := this.linkMap[linkItem]
		if linkExist {
			errRet = fmt.Errorf("mapping already exist")
			break
		}
		id, err := this.ormAddRightItem2RightSet(itemItem.Id, setItem.Id)
		if err != nil {
			errRet = err
			break
		}
		mapNode := UserRightSet2itemMap{
			Id:     id,
			ItemId: itemItem.Id,
			SetId:  setItem.Id,
		}
		this.linkMap[linkItem] = mapNode
		break
	}
	this.mutex.Unlock()
	return errRet
}

package db

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"sort"
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

type DefaultRightSets struct {
	Id         int64 `orm:"auto"`
	Name       string
	Dsc        string
	RightSetId int64
}

// ============User rights manager instance===========
type rightSet2ItemPair struct {
	setId  int64
	itemId int64
}

type DefaultRightSetMapNode struct {
	Name         string
	Dsc          string
	RightSetName string
}

type DbUserRightsManager struct {
	itemMap             map[string]*UserRightItem
	itemIdMap           map[int64]*UserRightItem
	setMap              map[string]*UserRightSet
	setIdMap            map[int64]*UserRightSet
	linkMap             map[rightSet2ItemPair]*UserRightSet2itemMap
	linkSetMap          map[int64]map[int64]*UserRightSet2itemMap
	defaultRightSetsMap map[string]*DefaultRightSets
	mutex               sync.Mutex
}

//Init 初始化
//装数据库数据
func (this *DbUserRightsManager) Init() {
	this.itemMap = make(map[string]*UserRightItem)
	this.itemIdMap = make(map[int64]*UserRightItem)
	this.setMap = make(map[string]*UserRightSet)
	this.setIdMap = make(map[int64]*UserRightSet)
	this.linkMap = make(map[rightSet2ItemPair]*UserRightSet2itemMap)
	this.linkSetMap = make(map[int64]map[int64]*UserRightSet2itemMap)
	this.defaultRightSetsMap = make(map[string]*DefaultRightSets)
	this.loadRightItems()
	this.loadRightSets()
	this.loadRightSet2ItemMapping()
	this.loadDefaultRightSets()
}

//仅供init使用
func (this *DbUserRightsManager) loadRightItems() {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	var rightItems []*UserRightItem
	o := orm.NewOrm()
	qs := o.QueryTable("user_right_item")
	num, err := qs.All(&rightItems)
	if num == 0 {
		return
	}
	if err != nil {
		panic(err)
	}
	for _, i := range rightItems {
		this.itemMap[i.Name] = i
		this.itemIdMap[i.Id] = i
	}
}

//仅供init使用
func (this *DbUserRightsManager) loadRightSets() {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	var rightSets []*UserRightSet
	o := orm.NewOrm()
	qs := o.QueryTable("user_right_set")
	num, err := qs.All(&rightSets)
	if num == 0 {
		return
	}
	if err != nil {
		panic(err)
	}
	for _, i := range rightSets {
		this.setMap[i.Name] = i
		this.setIdMap[i.Id] = i
	}
}

//仅供init使用
func (this *DbUserRightsManager) loadDefaultRightSets() {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	var rightSets []*DefaultRightSets
	o := orm.NewOrm()
	qs := o.QueryTable("default_right_sets")
	num, err := qs.All(&rightSets)
	if num == 0 {
		return
	}
	if err != nil {
		panic(err)
	}
	for _, i := range rightSets {
		this.defaultRightSetsMap[i.Name] = i
	}
}

func (this *DbUserRightsManager) checkRightItemExistanceWithIdWithoutLock(
	itemId int64) bool {
	_, b := this.itemIdMap[itemId]
	return b
}

func (this *DbUserRightsManager) checkRightSetExistanceWithIdWithoutLock(
	setId int64) bool {
	_, b := this.setIdMap[setId]
	return b
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
		this.linkMap[link] = i

		itemsMap := this.linkSetMap[i.SetId]
		if itemsMap == nil {
			itemsMap = make(map[int64]*UserRightSet2itemMap)
		}
		itemsMap[i.ItemId] = i
		this.linkSetMap[i.SetId] = itemsMap
	}
}

func (this *DbUserRightsManager) GetRightItems() []string {
	var keys []string
	this.mutex.Lock()
	for k := range this.itemMap {
		keys = append(keys, k)
	}
	this.mutex.Unlock()
	sort.Strings(keys)
	return keys
}

func (this *DbUserRightsManager) GetRightSets() []string {
	var keys []string
	this.mutex.Lock()
	for k := range this.setMap {
		keys = append(keys, k)
	}
	this.mutex.Unlock()
	sort.Strings(keys)
	return keys
}

func (this *DbUserRightsManager) HasRightItem(item string) bool {
	var b bool
	this.mutex.Lock()
	_, b = this.itemMap[item]
	this.mutex.Unlock()
	return b
}

func (this *DbUserRightsManager) GetRightSetRightItems(set string) []string {
	var items []string
	this.mutex.Lock()
	for {
		setData, bExist := this.setMap[set]
		if !bExist {
			break
		}
		itemMap, bExist := this.linkSetMap[setData.Id]
		if !bExist {
			break
		}
		for _, v := range itemMap {
			rightItemNode := this.itemIdMap[v.ItemId]
			if rightItemNode == nil {
				panic(fmt.Errorf("right item not exist while it is listed in link map"))
			}
			items = append(items, rightItemNode.Name)
		}
		break
	}
	this.mutex.Unlock()
	sort.Strings(items)
	return items
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
	var setItem *UserRightSet
	var itemItem *UserRightItem

	bRet = false
	this.mutex.Lock()
	for {
		setItem, setExist = this.setMap[set]
		itemItem, itemExist = this.itemMap[item]
		if !setExist || !itemExist {
			break
		}
		linkItem := rightSet2ItemPair{
			setId:  setItem.Id,
			itemId: itemItem.Id,
		}
		_, linkExist := this.linkMap[linkItem]
		if !linkExist {
			break
		}
		bRet = true
		break
	}
	this.mutex.Unlock()
	return bRet
}

func (this *DbUserRightsManager) GetDefaultRightSetsList() []*DefaultRightSetMapNode {
	var mapList []*DefaultRightSetMapNode
	this.mutex.Lock()
	for _, v := range this.defaultRightSetsMap {
		setNode := this.setIdMap[v.RightSetId]
		if setNode == nil {
			panic(fmt.Errorf("set %v not exist while it listed in default right sets map", v.RightSetId))
		}
		node := &DefaultRightSetMapNode{
			Name:         v.Name,
			Dsc:          v.Dsc,
			RightSetName: setNode.Name,
		}
		mapList = append(mapList, node)
	}
	this.mutex.Unlock()
	return mapList
}

//非线程安全
func (this *DbUserRightsManager) ormAddDeafultRightSet(
	name string, dsc string, setId int64) (int64, error) {
	var item DefaultRightSets
	item.Name = name
	item.Dsc = dsc
	item.RightSetId = setId

	o := orm.NewOrm()
	id, err := o.Insert(&item)
	if err != nil {
		logs.Error("[ORM]insert default user right set fail with error:", err)
	}
	return id, err
}

//非线程安全
func (this *DbUserRightsManager) ormUpdateDeafultRightSet(
	item *DefaultRightSets) error {
	o := orm.NewOrm()
	num, err := o.Update(item)
	if err != nil {
		logs.Error("[ORM]insert default user right set fail with error:", err)
	}
	if num == 0 {
		err = fmt.Errorf("[ORM]no rows impact when update DefaultRightSets with id=%v name=%s",
			item.Id, item.Name)
	}
	return err
}

func (this *DbUserRightsManager) AddDefaultRightSet(
	name string, dsc string, rightSetName string) error {
	var err error
	this.mutex.Lock()
	for {
		setData, bExist := this.setMap[rightSetName]
		if !bExist {
			err = fmt.Errorf("RightSet %s not exist", rightSetName)
			break
		}
		var id int64
		id, err = this.ormAddDeafultRightSet(name, dsc, setData.Id)
		if err != nil {
			break
		}
		this.defaultRightSetsMap[name] = &DefaultRightSets{
			Id:         id,
			Name:       name,
			Dsc:        dsc,
			RightSetId: setData.Id,
		}
		break
	}
	this.mutex.Unlock()
	return err
}

func (this *DbUserRightsManager) UpdateDefaultRightSet(
	name string, dsc string, rightSetName string) error {
	var err error
	this.mutex.Lock()
	for {
		mapNode, bExist := this.defaultRightSetsMap[name]
		if !bExist {
			err = fmt.Errorf("DefaultRightSet %s not exist", rightSetName)
			break
		}
		setData, bExist := this.setMap[rightSetName]
		if !bExist {
			err = fmt.Errorf("RightSet %s not exist", rightSetName)
			break
		}
		mapNode.Dsc = dsc
		mapNode.RightSetId = setData.Id
		err = this.ormUpdateDeafultRightSet(mapNode)
		if err != nil {
			break
		}
		this.defaultRightSetsMap[name] = mapNode
		break
	}
	this.mutex.Unlock()
	return err
}

//非线程安全
func (this *DbUserRightsManager) ormAddRightItem(
	name string, enabled bool, dsc string) (int64, error) {
	var item UserRightItem
	item.Name = name
	item.Discription = dsc
	item.Enable = enabled

	o := orm.NewOrm()
	id, err := o.Insert(&item)
	if err != nil {
		logs.Error("[ORM]insert user right item fail with error:%+v", err)
	}
	return id, err
}

func (this *DbUserRightsManager) AddRightItem(
	item string,
	enabled bool,
	dsc string) error {
	var errRet error
	this.mutex.Lock()
	for {
		_, b := this.itemMap[item]
		if b {
			errRet = fmt.Errorf("Right item [%s] already exists", item)
			break
		}
		id, err := this.ormAddRightItem(item, enabled, dsc)
		if err != nil {
			break
		}
		rightItem := UserRightItem{
			id, item, enabled, dsc,
		}
		this.itemMap[item] = &rightItem
		this.itemIdMap[id] = &rightItem
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
		errRet = this.ormUpdateRightItemRecord(r)
		if errRet == nil {
			this.itemMap[item] = r
			this.itemIdMap[r.Id] = r
		}
		break
	}
	this.mutex.Unlock()
	return errRet
}

func (this *DbUserRightsManager) RightItemEnabled(
	item string) (bool, error) {
	var errRet error
	enabled := false
	this.mutex.Lock()
	for {
		r, exist := this.itemMap[item]
		if !exist {
			errRet = fmt.Errorf("item %v not exist", item)
			break
		}
		enabled = r.Enable
		break
	}
	this.mutex.Unlock()
	return enabled, errRet
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
		this.setMap[set] = &setItem
		this.setIdMap[id] = &setItem
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

//非线程安全
func (this *DbUserRightsManager) ormDelRightItem2RightSet(
	linkId int64) error {

	mapNode := UserRightSet2itemMap{
		Id: linkId,
	}

	o := orm.NewOrm()
	_, err := o.Delete(&mapNode)
	if err != nil {
		logs.Error("[ORM]insert Set2ItemMap node fail with error:%+v", err)
	}
	return err
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
		this.linkMap[linkItem] = &mapNode
		setItemsMap := this.linkSetMap[setItem.Id]
		if setItemsMap == nil {
			setItemsMap = make(map[int64]*UserRightSet2itemMap)
		}
		setItemsMap[itemItem.Id] = &mapNode
		this.linkSetMap[setItem.Id] = setItemsMap
		break
	}
	this.mutex.Unlock()
	return errRet
}

func (this *DbUserRightsManager) DelRightItemFromRightSet(
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
		linkNode, linkExist := this.linkMap[linkItem]
		if !linkExist {
			errRet = fmt.Errorf("mapping not exist")
			break
		}
		err := this.ormDelRightItem2RightSet(linkNode.Id)
		if err != nil {
			errRet = err
			break
		}
		delete(this.linkMap, linkItem)
		setItemList := this.linkSetMap[setItem.Id]
		if setItemList == nil {
			panic(fmt.Errorf("Right item %s:%s listed in linkMap via, not exist in linkSetMap", set, item))
		}
		delete(setItemList, itemItem.Id)
		this.linkSetMap[setItem.Id] = setItemList

		break
	}
	this.mutex.Unlock()
	return errRet
}

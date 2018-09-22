package models

import (
	"atnero.com/blog/models/db"
)

type UserRightsManagerIf interface {
	Init()
	GetRightItems() []string
	GetRightSets() []string
	GetRightSetId(set string) (int64, error)
	HasRightItem(item string) bool
	HasRightSet(set string) bool
	RightSetHasRightItem(set string, item string) bool
	AddRightItem(item string, enabled bool, dsc string) error
	EnableRightItem(item string, enable bool) error
	RightItemEnabled(item string) (bool, error)
	GetRightItemDiscription(item string) (string, bool)
	AddRightSet(set string, dsc string) error
	GetRightSetNameById(id int64) (string, error)
	GetRightSetDiscription(set string) (string, bool)
	AddRightItem2RightSet(set string, item string) error
	DelRightItemFromRightSet(set string, item string) error
	GetRightSetRightItems(set string) []string
	GetDefaultRightSetsList() []db.DefaultRightSetMapNode
	AddDefaultRightSet(name string, dsc string, rightSetName string) error
	UpdateDefaultRightSet(name string, dsc string, rightSetName string) error
	GetDefaultRightSetName(name string) (string, error)
	HasDefaultRightSet(name string) bool
}

var userRightsManagerInst UserRightsManagerIf

func init() {
	userRightsManagerInst = new(db.DbUserRightsManager)
	userRightsManagerInst.Init()
}

func UserRightsMngInst() UserRightsManagerIf {
	return userRightsManagerInst
}

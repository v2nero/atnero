package models

import (
	"atnero.com/blog/models/db"
)

type UserRightsManagerIf interface {
	Init()
	GetRightItems() []string
	GetRightSets() []string
	HasRightItem(item string) bool
	HasRightSet(set string) bool
	RightSetHasRightItem(set string, item string) bool
	AddRightItem(item string, enabled bool, dsc string) error
	EnableRightItem(item string, enable bool) error
	RightItemEnabled(item string) (bool, error)
	GetRightItemDiscription(item string) (string, bool)
	AddRightSet(set string, dsc string) error
	GetRightSetDiscription(set string) (string, bool)
	AddRightItem2RightSet(set string, item string) error
}

var userRightsManagerInst UserRightsManagerIf

func init() {
	userRightsManagerInst = new(db.DbUserRightsManager)
	userRightsManagerInst.Init()
}

func UserRightsMngInst() UserRightsManagerIf {
	return userRightsManagerInst
}

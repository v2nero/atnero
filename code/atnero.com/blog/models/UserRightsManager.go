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
	AddRightItem(item string, dsc string) error
	AddRightSet(set string, dsc string) error
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

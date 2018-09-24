package db

import (
	"time"
)

type Users struct {
	Id       int64 `orm:"auto"`
	Name     string
	Password string `orm:"column(pwd)"`
	Email    string
	Rightset int64
	RegTime  time.Time `orm:"auto_now_add;type(datetime)"`
	LastTime time.Time `orm:"auto_now_add;type(datetime)"`
	FailTime time.Time `orm:"auto_now_add;type(datetime)"`
}
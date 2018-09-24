package db

import ()
import "time"

type InvitationCode struct {
	Id         int64 `orm:"auto"`
	Code       string
	Used       bool
	CreateTime time.Time `orm:"auto_now;type(datatime)"`
	ExpireTime time.Time `orm:"type(datatime)"`
}

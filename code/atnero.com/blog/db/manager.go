package db

type DbVersion struct {
	Id      int64 `orm:"auto"`
	Version string
}

type DbManager struct {
	Id     int64 `orm:"auto"`
	Enable bool
}

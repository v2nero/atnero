package db

type Version struct {
	Id      int64 `orm:"auto"`
	Version string
}

/*
type BgManagerEnable struct {
	Id     int64 `orm:"auto"`
	Enable bool
}
*/

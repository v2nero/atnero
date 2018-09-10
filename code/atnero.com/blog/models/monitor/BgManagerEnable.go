package monitor

import (
	"atnero.com/blog/models/db"
	"fmt"
	"net"
)

type CmdBgMngEnable struct {
}

func (this *CmdBgMngEnable) Name() string {
	return "BgMngEnable"
}

func (this *CmdBgMngEnable) Handle(conn net.Conn, args ...string) bool {
	pwd := db.DbMgrInst().EnableBgManager()
	fmt.Fprintln(conn, "请使用验证码", pwd, "从/manager/userright登陆")
	return false
}

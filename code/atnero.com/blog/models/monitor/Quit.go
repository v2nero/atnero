package monitor

import (
	"net"
)

type CmdQuit struct {
}

func (this *CmdQuit) Name() string {
	return "Quit"
}

func (this *CmdQuit) Handle(conn net.Conn, args ...string) bool {
	return true
}

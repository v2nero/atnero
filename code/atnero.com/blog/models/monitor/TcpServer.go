package monitor

import (
	"bufio"
	"fmt"
	"github.com/astaxie/beego/logs"
	"net"
	"strings"
)

var connectLimits chan int

type CmdHandleIf interface {
	Name() string
	Handle(conn net.Conn, args ...string) bool //返回值用于判断是否要断开连接
}

var handlers map[string]CmdHandleIf

func getArgs(cmdLine string) (params []string, err error) {
	scanner := bufio.NewScanner(strings.NewReader(cmdLine))
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		params = append(params, scanner.Text())
	}
	if scannerErr := scanner.Err(); scannerErr != nil {
		err = fmt.Errorf("fail to parse cmd line %s, err=%v", cmdLine, scannerErr)
	}
	return
}

func cmdHandleLoop(conn net.Conn) {
	input := bufio.NewScanner(conn)
	for input.Scan() {
		cmdLine := input.Text()
		args, err := getArgs(cmdLine)
		if err != nil {
			logs.Error(err)
		} else if len(args) > 0 {
			h, ok := handlers[args[0]]
			if !ok {
				logs.Error("[Monitor] cmd ", args[0], " not exist")
			} else {
				bRet := h.Handle(conn, args...)
				if bRet {
					logs.Info("[Monitor] cmd ", args[0], " query to quit")
					break
				}
			}
		}
	}
	<-connectLimits
}

func listenRoutine(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			logs.Error(err)
			continue
		}
		connectLimits <- 1
		go cmdHandleLoop(conn)
	}
}

func InitServer() error {
	connectLimits = make(chan int, 1)
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		logs.Error("[Monitor]", err)
		return err
	}
	go listenRoutine(listener)
	return nil
}

func registCmdHandler(h CmdHandleIf) {
	name := h.Name()
	_, ok := handlers[name]
	if ok {
		logs.Error("[Monitor] cmd ", name, " already exist")
		return
	}
	handlers[name] = h
}

func init() {
	handlers = make(map[string]CmdHandleIf)

	//注册命令参数
	registCmdHandler(new(CmdQuit))
	registCmdHandler(new(CmdBgMngEnable))
}

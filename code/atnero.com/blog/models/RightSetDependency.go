package models

import (
	"fmt"
	"github.com/astaxie/beego"
)

type myRightSetDependencyNode struct {
	name string //名字
	src  string //依赖源
	dsc  string //描述
}

var myRightSetDependencies map[string]myRightSetDependencyNode

func AddDependencyRightSet(name string, src string, dsc string) {
	if len(name) == 0 {
		panic(fmt.Errorf("invalid dependency name"))
	}
	_, bExist := myRightSetDependencies[name]
	if bExist {
		return
	}
	myRightSetDependencies[name] = myRightSetDependencyNode{
		name: name,
		src:  src,
		dsc:  dsc,
	}
}

func CheckRightSetDependencies() bool {
	meet := true
	for k := range myRightSetDependencies {
		if !UserRightsMngInst().HasRightSet(k) {
			meet = false
			beego.Error("Missing RightSet ", k)
		}
	}
	return meet
}

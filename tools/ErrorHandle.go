package tools

import (
	"fmt"
	"runtime"
)

func ErrorLogging(err error) {
	pc, file, line, ok := runtime.Caller(1)
	pcName := runtime.FuncForPC(pc).Name() //获取函数名
	fmt.Println("ERROR: ", err)
	fmt.Println(fmt.Sprintf("%v   %s   %d   %t   %s", pc, file, line, ok, pcName))
	return
}

package main

import (
	"fmt"
	"runtime"
)

func test() {

}

func main() {
	funcName, file, line, ok := runtime.Caller(0)
	test()
	if ok {
		fmt.Println("file:", file, " func:", runtime.FuncForPC(funcName).Name(), " line:", line)
	}
}

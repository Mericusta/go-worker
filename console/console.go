// +build windows

package console

import (
	"syscall"

	"github.com/go-worker/ui"
)

var dllKernel32 *syscall.DLL
var dllKeyKernal32 = "kernel32.dll"

var procReadConsoleInputW *syscall.Proc
var procKeyreadConsoleInputW = "ReadConsoleInputW"

// 目前 DLL 库只用于监听键盘方向键回溯上次输入的指令
// 对整体逻辑无任何影响，所以即使出错了也不需要退出程序

func init() {
	var loadDLLError error
	dllKernel32, loadDLLError = syscall.LoadDLL(dllKeyKernal32)
	if loadDLLError != nil {
		ui.OutputErrorInfo(ui.CONSOLELoadDLLError, dllKeyKernal32, loadDLLError)
	}

	var findProcError error
	procReadConsoleInputW, findProcError = dllKernel32.FindProc(procKeyreadConsoleInputW)
	if findProcError != nil {
		ui.OutputErrorInfo(ui.CONSOLEFindProcError, procKeyreadConsoleInputW, dllKeyKernal32, findProcError)
	}
}

func Run() {

}

func readConsoleInput() {

}

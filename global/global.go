package global

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/go-worker/fsm"
)

var OperationSystem string

var FsmState fsm.FSMState

var SyntaxFileSuffixMap map[string]string

var GoPathSrc string

var WorkerRootPath string

func init() {
	OperationSystem = runtime.GOOS
	SyntaxFileSuffixMap = make(map[string]string)
	SyntaxFileSuffixMap[SyntaxGo] = "go"
	SyntaxFileSuffixMap[SyntaxCSV] = "csv"
	if OperationSystem == OSWindows {
		GoPathSrc = fmt.Sprintf("%v\\%v", os.Getenv("GOPATH"), "src\\")
	} else if OperationSystem == OSLinux {
		GoPathSrc = fmt.Sprintf("%v/%v", os.Getenv("GOPATH"), "src/")
	}
	WorkerRootPath = filepath.Join(GoPathSrc, "go-worker")
}

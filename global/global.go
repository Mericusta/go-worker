package global

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/go-worker/fsm"
)

var FsmState fsm.FSMState

var SyntaxFileSuffixMap map[string]string

var GoPathSrc string

var WorkerRootPath string

func init() {
	SyntaxFileSuffixMap = make(map[string]string)
	SyntaxFileSuffixMap[SyntaxGo] = "go"
	SyntaxFileSuffixMap[SyntaxCSV] = "csv"
	GoPathSrc = strings.Replace(fmt.Sprintf("%v\\%v", os.Getenv("GOPATH"), "src\\"), "\\", "/", -1)
	WorkerRootPath = path.Join(GoPathSrc, "go-worker")
}

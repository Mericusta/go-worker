package config

import (
	"path/filepath"

	"github.com/go-worker/global"
	"github.com/go-worker/utility"
)

type Config struct {
	ProjectPath        string
	ProjectSyntax      string
	ConvertCsvHeadLine int
	ConvertCsvTypeLine int
	ConvertCsvSplitter string
	LogFilePath        string
}

var WorkerConfig *Config

func init() {
	initWorkerConfig()
}

func initWorkerConfig() {
	WorkerConfig = &Config{
		ProjectSyntax:      global.SyntaxGo,
		ConvertCsvHeadLine: 2,
		ConvertCsvTypeLine: 1,
		ConvertCsvSplitter: "\t",
	}
	WorkerConfig.LogFilePath = filepath.Join(GetCurrentProjectPath(), "worker.log")
}

// GetCurrentProjectPath 获取当前绑定的项目
func GetCurrentProjectPath() string {
	path := "."
	projectPath := WorkerConfig.ProjectPath
	if projectPath != "" && utility.IsExist(projectPath) {
		path = projectPath
	}
	return path
}

// GetCurrentSyntaxFileSuffix 获取当前绑定语法的文件扩展名
func GetCurrentSyntaxFileSuffix() string {
	fileType := ""
	projectSyntax := global.SyntaxFileSuffixMap[WorkerConfig.ProjectSyntax]
	if projectSyntax != "" {
		fileType = projectSyntax
	}
	return fileType
}

// GetSpecificSyntaxFileSuffix 获取指定语法的文件扩展名
func GetSpecificSyntaxFileSuffix(syntax string) string {
	fileType := ""
	projectSyntax := global.SyntaxFileSuffixMap[syntax]
	if projectSyntax != "" {
		fileType = projectSyntax
	}
	return fileType
}

package commands

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/go-worker/config"
	"github.com/go-worker/global"
	"github.com/go-worker/ui"
	"github.com/go-worker/utility"
)

type Convert struct {
	*CommandStruct
	Params *convertParam
}

func (command *Convert) Execute() error {
	// 解析指令的选项和参数
	parseCommandParamsError := command.parseCommandParams()
	if parseCommandParamsError != nil {
		return parseCommandParamsError
	}

	// 获取全局配置
	projectPath := config.GetCurrentProjectPath()
	fileType := config.GetSpecificSyntaxFileSuffix(global.SyntaxCSV)
	if fileType == "" {
		ui.OutputWarnInfo(ui.CommonWarn1)
	}

	utility.TestOutput("projectPath = %v, fileType = %v", projectPath, fileType)
	utility.TestOutput("command.Params = %+v", command.Params)

	filePath := fmt.Sprintf("%v/%v.%v", projectPath, command.Params.sourceValue, fileType)
	if command.Params.sourceParentValue != "" {
		filePath = fmt.Sprintf("%v/%v/%v.%v", projectPath, command.Params.sourceParentValue, command.Params.sourceValue, fileType)
	}

	file, inputError := os.Open(filePath)
	defer file.Close()
	if inputError != nil || file == nil {
		return fmt.Errorf("open file %v error, file is nil or %v", filePath, inputError.Error())
	}

	fileContent := ""
	var convertError error
	switch command.Params.sourceType {
	case "csv":
		fileContent, convertError = convertCsvToStruct(command.Params.sourceValue, file)
	default:
		ui.OutputNoteInfo(ui.CommonNote1)
		return nil
	}

	if convertError != nil {
		return convertError
	}

	utility.TestOutput("fileContent = %v", fileContent)

	var toWriteFile *os.File
	defer func() {
		if toWriteFile != nil {
			toWriteFile.Close()
		}
	}()
	toWriteFileType := config.GetCurrentSyntaxFileSuffix()
	toWriteFilePath := fmt.Sprintf("%v.%v", command.Params.targetValue, toWriteFileType)
	switch command.Params.targetOption {
	case "create", "":
		if command.Params.targetOption == "" {
			toWriteFilePath = fmt.Sprintf("%v.%v", utility.Convert2CamelStyle(command.Params.sourceValue, false), toWriteFileType)
			utility.TestOutput("toWriteFilePath = %v", toWriteFilePath)
		}
		if utility.IsExist(toWriteFilePath) {
			var openFileError error
			toWriteFile, openFileError = os.Open(toWriteFilePath)
			if openFileError != nil {
				utility.TestOutput("error 1")
				return openFileError
			}
		} else {
			var createFileError error
			toWriteFile, createFileError = utility.CreateFile(toWriteFilePath)
			if createFileError != nil {
				utility.TestOutput("error 2")
				return createFileError
			}
		}
	case "append":
		if utility.IsExist(toWriteFilePath) {
			var openFileError error
			toWriteFile, openFileError = os.Open(toWriteFilePath)
			if openFileError != nil {
				utility.TestOutput("error 3")
				return openFileError
			}
		}
	default:
		ui.OutputNoteInfo(ui.CommonNote1)
		return nil
	}

	if toWriteFile == nil {
		utility.TestOutput("error 4")
		return fmt.Errorf("convert to write file is nil")
	}

	_, writeError := toWriteFile.WriteString(fileContent)
	if writeError != nil {
		utility.TestOutput("error 5")
		return writeError
	}

	return nil
}

const (
	convertCommandIndex     = 0
	convertSourceTypeIndex  = 1
	convertSourceValueIndex = 2
)

type convertParam struct {
	sourceType        string
	sourceValue       string
	sourceParentValue string
	targetOption      string
	targetValue       string
}

func (command *Convert) parseCommandParams() error {
	inputStringList := strings.Split(command.CommandStruct.InputString, " ")
	if convertSourceValueIndex >= len(inputStringList) {
		return fmt.Errorf(ui.CommonError1)
	}
	_, inputStringList = utility.SlicePop(inputStringList)
	pop := func() string {
		var element string
		element, inputStringList = utility.SlicePop(inputStringList)
		return element
	}
	command.Params = &convertParam{
		sourceType:  pop(),
		sourceValue: pop(),
		sourceParentValue: func() string {
			if len(inputStringList) > 1 {
				element, list := utility.SlicePop(inputStringList)
				if element == "parent" {
					inputStringList = list
					return pop()
				}
			}
			return ""
		}(),
		targetOption: func() string {
			if len(inputStringList) > 1 {
				return pop()
			}
			return ""
		}(),
		targetValue: func() string {
			if len(inputStringList) > 0 {
				return pop()
			}
			return ""
		}(),
	}
	return nil
}

func convertCsvToStruct(fileName string, file *os.File) (string, error) {
	fileReader := bufio.NewReader(file)

	headStringList := make([]string, 0)
	typeStringList := make([]string, 0)

	line := 0
	for {
		lineString, readerError := fileReader.ReadString('\n')
		if readerError != nil && readerError != io.EOF {
			ui.OutputErrorInfo("read line[%v] error: %v", line, readerError)
			continue
		}
		if readerError == io.EOF {
			break
		}

		lineTrimSpace := strings.TrimSpace(lineString)

		if line == config.WorkerConfig.ConvertCsvHeadLine {
			headStringList = strings.Split(lineTrimSpace, config.WorkerConfig.ConvertCsvSplitter)
		} else if line == config.WorkerConfig.ConvertCsvTypeLine {
			typeStringList = strings.Split(lineTrimSpace, config.WorkerConfig.ConvertCsvSplitter)
		} else if line > config.WorkerConfig.ConvertCsvHeadLine && line > config.WorkerConfig.ConvertCsvTypeLine {
			break
		}

		line++
	}

	if len(headStringList) != len(typeStringList) {
		return "", fmt.Errorf("parse csv head list length[%v] not equal type list length[%v]", len(headStringList), len(typeStringList))
	}

	memberContent := ""
	for index := 0; index != len(headStringList); index++ {
		memberCommentFormatString := utility.ConvertTemplate2Format(ui.GoMemberCommentByCSV)
		memberCommentString := fmt.Sprintf(memberCommentFormatString, headStringList[index])
		memberFormatString := utility.ConvertTemplate2Format(ui.GoStructMemberTemplate)
		memberString := fmt.Sprintf(memberFormatString, utility.Convert2CamelStyle(headStringList[index], true), typeStringList[index], memberCommentString)
		if memberContent == "" {
			memberContent = memberString
		} else {
			memberContent = fmt.Sprintf("%v\n\t%v", memberContent, memberString)
		}
	}

	structFormatString := utility.ConvertTemplate2Format(ui.GoStructTemplate)
	structContent := fmt.Sprintf(structFormatString, utility.Convert2CamelStyle(fileName, true), memberContent)

	return structContent, nil
}

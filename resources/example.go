package regexps

// GetRegexpByTemplateEnum 根据模板枚举获得其原子解析式
func GetRegexpByTemplateEnum(templateEnum global.TemplateEnum) *regexp.Regexp {
	templateExpression, hasTemplateExpression := template.TemplateExpressionMap[templateEnum]
	if !hasTemplateExpression {
		ui.OutputErrorInfo("template expression[%v] does not regist expression", templateEnum)
		return nil
	}
	templateCommonKeywordRegexp, hasTemplateCommonKeywordRegexp := AtomicExpressionEnumRegexpMap[global.AETemplateCommonKeyword]
	if !hasTemplateCommonKeywordRegexp {
		ui.OutputErrorInfo("template common keyword regexp does not exist")
		return nil
	}
	parsedCommand := parseTemplateExpression(templateCommonKeywordRegexp, templateExpression)
	if parsedCommand == "" {
		ui.OutputWarnInfo("parse template expression[%v] but get empty", templateExpression)
		return nil
	}
	templateExpressionRegexp := regexp.MustCompile(string(parsedCommand))
	if templateExpressionRegexp == nil {
		ui.OutputWarnInfo("complie template expression[%v], but get nil", parsedCommand)
	}
	return templateExpressionRegexp
}

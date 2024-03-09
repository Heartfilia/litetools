package litejson

import "github.com/Heartfilia/litetools/litejson/jsonPath"

func TryGet(jsonString, pathNode string) (jsonPath.Result, error) {
	// 本来不想这样子返回两个参数的 但是没办法 只能这样子返回了
	result := jsonPath.Result{}
	newRule := jsonPath.SplitRule(pathNode)
	jsonPath.JudgeAndExtractEachRule(jsonString, newRule, &result)
	return result, result.Error
}

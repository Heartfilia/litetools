package litejson

import "github.com/Heartfilia/litetools/litejson/jsonPath"

func TryGet(jsonString, pathNode string) jsonPath.Result {
	// 还是改成只返回一个值的好
	result := jsonPath.Result{}
	newRule := jsonPath.SplitRule(pathNode)
	jsonPath.JudgeAndExtractEachRule(jsonString, newRule, &result)
	return result
}

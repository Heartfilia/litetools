package litejson

import "github.com/Heartfilia/litetools/litejson/jsonPath"

func TryGet(jsonString, rulePath string) jsonPath.Result {
	// 设计成这样子就是为了避免获取到error  所以我把error写进了 Result
	result := jsonPath.Result{}
	newRule := jsonPath.SplitRule(rulePath)
	jsonPath.JudgeAndExtractEachRule(jsonString, newRule, &result)
	return result
}

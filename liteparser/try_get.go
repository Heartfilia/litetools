package liteparser

import "github.com/Heartfilia/litetools/liteparser/jsonPath"

func TryGet(jsonString, rulePath string) (*jsonPath.Result, error) {
	newRule := jsonPath.SplitRule(rulePath)
	jsonPath.JudgeAndExtractEachRule(jsonString, newRule)
	return &jsonPath.Result{}, nil
}

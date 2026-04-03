package jsonPath

import (
	"strings"
	"testing"
)

func TestReplaceToRoundTrip(t *testing.T) {
	input := `user\.name[0]\|hello\ world`

	encoded := replaceTo(input, 1)
	decoded := replaceTo(encoded, 0)

	if decoded != input {
		t.Fatalf("round-trip mismatch: got %q want %q", decoded, input)
	}
}

func TestJudgeAndExtractEachRuleNegativeIndex(t *testing.T) {
	jsonString := `{"items":[{"id":1},{"id":2},{"id":3}]}`
	var result Result

	JudgeAndExtractEachRule(jsonString, []string{"items[-1].id"}, &result)

	if result.Err != nil {
		t.Fatalf("unexpected error: %v", result.Err)
	}
	if got := result.Int(); got != 3 {
		t.Fatalf("unexpected value: got %d want %d", got, 3)
	}
}

func TestJudgeAndExtractEachRuleResetsResultState(t *testing.T) {
	jsonString := `{"items":[{"id":1}]}`
	var result Result

	JudgeAndExtractEachRule(jsonString, []string{"items[0].id"}, &result)
	if result.Err != nil || result.Int() != 1 {
		t.Fatalf("expected first extraction to succeed, got value=%v err=%v", result.Value(), result.Err)
	}

	JudgeAndExtractEachRule(jsonString, []string{"missing"}, &result)
	if result.Err != nil {
		t.Fatalf("expected missing terminal key to return nil without error, got %v", result.Err)
	}
	if result.Value() != nil {
		t.Fatalf("expected result to be reset to nil, got %v", result.Value())
	}
}

func TestJudgeAndExtractEachRuleInvalidSegmentReturnsError(t *testing.T) {
	jsonString := `{"items":[{"id":1}]}`
	var result Result

	JudgeAndExtractEachRule(jsonString, []string{"items]"}, &result)

	if result.Err == nil {
		t.Fatal("expected invalid rule segment to return error")
	}
	if !strings.Contains(result.Err.Error(), "unsupported rule segment") {
		t.Fatalf("unexpected error: %v", result.Err)
	}
}

func TestJudgeAndExtractEachRuleWrongNodeTypeReturnsError(t *testing.T) {
	jsonString := `{"items":{"id":1}}`
	var result Result

	JudgeAndExtractEachRule(jsonString, []string{"items[0]"}, &result)

	if result.Err == nil {
		t.Fatal("expected array access on object to return error")
	}
	if !strings.Contains(result.Err.Error(), "not an array") {
		t.Fatalf("unexpected error: %v", result.Err)
	}
}

func TestJudgeAndExtractEachRuleEscapedDotKey(t *testing.T) {
	jsonString := `{"user.name":{"id":7}}`
	var result Result

	JudgeAndExtractEachRule(jsonString, []string{`user\.name.id`}, &result)

	if result.Err != nil {
		t.Fatalf("unexpected error: %v", result.Err)
	}
	if got := result.Int(); got != 7 {
		t.Fatalf("unexpected value: got %d want %d", got, 7)
	}
}

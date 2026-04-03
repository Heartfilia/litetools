package litereq

import (
	"io"
	"net/url"
	"testing"
)

func TestBuilderParamsPreservesDuplicatesAndEncodedValues(t *testing.T) {
	u, err := Build("https://example.com/api?base=1").
		Param("tag", "first").
		Param("tag", "second").
		Params("name=alice%20bob&empty&tag=third").
		GetUrl()
	if err != nil {
		t.Fatalf("GetUrl() error = %v", err)
	}

	got := u.RawQuery
	want := "base=1&empty=&name=alice%20bob&tag=first&tag=second&tag=third"
	if got != want {
		t.Fatalf("RawQuery = %q, want %q", got, want)
	}
}

func TestBodyDataSupportsDuplicateAndSliceValues(t *testing.T) {
	getter := bodyData(map[string]any{
		"tag":   []string{"first", "second"},
		"name":  "alice bob",
		"empty": nil,
		"mix":   []any{"x", 2},
	})

	body, err := getter()
	if err != nil {
		t.Fatalf("body getter error = %v", err)
	}
	defer body.Close()

	data, err := io.ReadAll(body)
	if err != nil {
		t.Fatalf("ReadAll() error = %v", err)
	}

	values, err := url.ParseQuery(string(data))
	if err != nil {
		t.Fatalf("ParseQuery() error = %v", err)
	}

	assertValues(t, values, "tag", []string{"first", "second"})
	assertValues(t, values, "name", []string{"alice bob"})
	assertValues(t, values, "empty", []string{""})
	assertValues(t, values, "mix", []string{"x", "2"})
}

func TestParseParamStringFallsBackForInvalidEscapes(t *testing.T) {
	values, err := parseParamString("note=100%&empty")
	if err != nil {
		t.Fatalf("parseParamString() error = %v", err)
	}

	assertValues(t, values, "note", []string{"100%"})
	assertValues(t, values, "empty", []string{""})
}

func assertValues(t *testing.T, values url.Values, key string, want []string) {
	t.Helper()

	got := values[key]
	if len(got) != len(want) {
		t.Fatalf("%s values len = %d, want %d (got %v)", key, len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("%s value[%d] = %q, want %q", key, i, got[i], want[i])
		}
	}
}

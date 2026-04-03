package litereq

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

func normalizeParams(v any) (url.Values, error) {
	values := make(url.Values)
	if v == nil {
		return values, nil
	}

	switch data := v.(type) {
	case url.Values:
		for key, items := range data {
			values[key] = append([]string(nil), items...)
		}
		return values, nil
	case map[string]string:
		for key, value := range data {
			values.Add(key, value)
		}
		return values, nil
	case map[string][]string:
		for key, items := range data {
			values[key] = append([]string(nil), items...)
		}
		return values, nil
	case map[string]any:
		for key, value := range data {
			for _, item := range flattenParamValue(value) {
				values.Add(key, item)
			}
		}
		return values, nil
	case string:
		return parseParamString(data)
	default:
		return nil, fmt.Errorf("unsupported param type %T", v)
	}
}

func parseParamString(raw string) (url.Values, error) {
	if raw == "" {
		return make(url.Values), nil
	}

	values, err := url.ParseQuery(raw)
	if err == nil {
		return values, nil
	}

	fallback := make(url.Values)
	for _, part := range strings.Split(raw, "&") {
		if part == "" {
			continue
		}
		key, value, _ := strings.Cut(part, "=")
		fallback.Add(key, value)
	}
	return fallback, nil
}

func flattenParamValue(v any) []string {
	if v == nil {
		return []string{""}
	}

	switch data := v.(type) {
	case []string:
		return append([]string(nil), data...)
	case []any:
		out := make([]string, 0, len(data))
		for _, item := range data {
			out = append(out, flattenParamValue(item)...)
		}
		return out
	case fmt.Stringer:
		return []string{data.String()}
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		if rv.Type().Elem().Kind() == reflect.Uint8 {
			return []string{string(rv.Bytes())}
		}
		out := make([]string, 0, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			out = append(out, flattenParamValue(rv.Index(i).Interface())...)
		}
		return out
	}

	return []string{fmt.Sprint(v)}
}

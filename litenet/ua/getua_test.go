package ua

import (
	"strings"
	"testing"

	"github.com/Heartfilia/litetools/utils/types"
)

func TestOptionsIOSUsesSafariFamily(t *testing.T) {
	localCacheMu.Lock()
	localCache = defaultSetting
	localCacheMu.Unlock()

	uaString := Options("ios")
	if !strings.Contains(uaString, "iPhone") {
		t.Fatalf("expected ios UA, got %q", uaString)
	}
	if !strings.Contains(uaString, "Safari/") {
		t.Fatalf("expected safari UA, got %q", uaString)
	}
	if strings.Contains(uaString, "Edg/") || strings.Contains(uaString, "OPR/") {
		t.Fatalf("unexpected non-ios browser marker in UA: %q", uaString)
	}
}

func TestOptionsLinuxDoesNotUseSafari(t *testing.T) {
	localCacheMu.Lock()
	localCache = defaultSetting
	localCacheMu.Unlock()

	uaString := CombineString("linux", browserForPlatform("linux", "safari"))
	if strings.Contains(uaString, "Mac OS X") {
		t.Fatalf("expected linux UA, got %q", uaString)
	}
}

func TestConfigFromCacheReturnsLocalCache(t *testing.T) {
	want := types.ConfigJson{
		Chromium: []string{"1.0.0.0"},
		Firefox:  []string{"1"},
		Safari:   []string{"1.0"},
	}

	localCacheMu.Lock()
	localCache = want
	localCacheMu.Unlock()

	got := configFromCache()
	if len(got.Chromium) != 1 || got.Chromium[0] != "1.0.0.0" {
		t.Fatalf("unexpected cache result: %+v", got)
	}
}

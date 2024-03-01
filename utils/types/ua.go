package types

type ConfigJson struct {
	Chromium []string `json:"chromium"`
	Firefox  []string `json:"firefox"`
	Safari   []string `json:"safari"`
}

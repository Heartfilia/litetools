package types

type ConfigJson struct {
	Chromium []string `json:"chromium"`
	Firefox  []string `json:"firefox"`
	Safari   []string `json:"safari"`
}

func (c *ConfigJson) IsEmpty() bool {
	return len(c.Chromium) == 0 && len(c.Firefox) == 0 && len(c.Safari) == 0
}

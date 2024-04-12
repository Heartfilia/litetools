package litepage

type Browser struct {
	Addr      string
	BrowserId int
	Page      int
}

func CreateBrowser(browser *Browser) *Browser {

	return browser
}

func (b *Browser) RunCDP(cmd string, cmdArgs map[string]string) {
	if cmdArgs != nil {

	}
}

package request

import (
	"io"
	"log"
	"net/http"
)

// 因为我这个项目只会用到get请求 所以我这里直接写死

func DoGet(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)
	item, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	return string(item)
}

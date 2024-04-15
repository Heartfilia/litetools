package litereq

import (
	"github.com/Heartfilia/litetools/litetime"
	"sync"
)

type Counter struct {
	// 用于统计计数的操作
	startTime int            // 开始时间 --> ms
	timeMs    int            // 程序运行的毫秒数
	crawlPage int            // 请求的页面数-含重试
	okRequest int            // 成功的请求
	pages     map[string]int // {统计每个域名下面的请求次数}
	lock      *sync.RWMutex
}

func NewCounter() *Counter {
	return &Counter{
		startTime: litetime.Time(nil).Int(),
		timeMs:    0,
		crawlPage: 0,
		okRequest: 0,
		lock:      &sync.RWMutex{},
	}
}

func (c *Counter) Default() {
	// 重置数据
	c.lock.Lock()
	defer c.lock.Unlock()

	c.startTime = litetime.Time(nil).Int()
	c.timeMs = 0
	c.crawlPage = 0
	c.okRequest = 0
}

func (c *Counter) ReqPage() {

}

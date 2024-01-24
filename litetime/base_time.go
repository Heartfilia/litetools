package litetime

type Time struct {
	Goal   interface{} // 基础数据类型 不传入 默认进行的是时间戳获取
	Fmt    interface{} // 格式化样式 不传入 默认不操作
	Unit   string      // 时间样式 s为秒 ms为毫秒
	Cursor interface{} // 游标 默认为0
}

func (t *Time) init() {
	if t.Unit == "" {
		t.Unit = "s"
	}
	if t.Cursor == nil {
		t.Cursor = 0
	}
}

func (t *Time) GetTime() interface{} {
	t.init()

	return nil
}

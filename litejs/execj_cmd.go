package litejs

import (
	"github.com/Heartfilia/litetools/litestr"
	"log"
	"os/exec"
)

// 只弄这个操作... 通过命令行的方式执行js
// 所以 对应的js内部需要完成调用 这里只捕获终端最后能打印的内容 故 需要得到 console.log的输出结果
// node js_path.js xxxx  xxxx

type CmdNode struct {
	Node    string // node 的执行路径，默认就是 node, 可以传入指定的路径的node
	JsPath  string // 对应的js的路径
	Verbose bool   // 是否需要打印提示 - 默认不打印
}

// Call
// 调用原生的 node 来获取结果
func (c *CmdNode) Call(args ...string) []byte {
	if c.Node == "" {
		c.Node = "node"
	}

	cmd := exec.Command(c.Node, append([]string{c.JsPath}, args...)...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("%s Error executing JavaScript file: %v\n", litestr.E(), err)
		return nil
	}
	if len(output) == 0 && c.Verbose == true {
		log.Printf(`%s
%s
+------------------------------------------------------------------------------------------+
|   (function(args){console.log(%s(...args));})(process.argv.slice(2))   |
+------------------------------------------------------------------------------------------+
%s：如果js返回的是一个对象最好是<JSON.stringify()>处理一次,要不然golang的json包会报错
`, litestr.W(), litestr.ColorString("将下面的内容放到你的js代码最后一行即可 调用的函数 换成你自己的", "red"), litestr.ColorString("这里换成你的函数名称", "green"),
			litestr.ColorString("特别注意", "yellow"))
	}
	return output
}

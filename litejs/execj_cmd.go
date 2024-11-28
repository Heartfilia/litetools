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

func (c *CmdNode) Call(args ...string) string {
	if c.Node == "" {
		c.Node = "node"
	}

	cmd := exec.Command(c.Node, append([]string{c.JsPath}, args...)...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("%s Error executing JavaScript file: %v\n", litestr.E(), err)
		return ""
	}
	res := string(output)
	if res == "" && c.Verbose == true {
		log.Printf(`%s 可能js没有按照要求进行适配：这个脚本是获取 命令行执行 node 后console.log得到的数据,所以需要 js那边调用，示例如下:
实际是 >>> node xxx.js 参数1 参数2 ....   无参数可不加
将下面的内容放到你的js代码最后一块即可 调用的函数 换成你自己的
----------------------------------------------------------------------------
%s %s() {
    %s args = process.argv.slice(2); // <<< 如果有参数的话从这里解析得到
    %s.%s(这里换成你的函数(...args));
}

%s()   // 这里就是调用方法
----------------------------------------------------------------------------
`, litestr.W(),
			litestr.ColorString("function", "yellow"), litestr.ColorString("processArguments", "blue"),
			litestr.ColorString("const", "yellow"), litestr.ColorString("console", "purple"),
			litestr.ColorString("log", "blue"), litestr.ColorString("processArguments", "blue"),
		)
	}
	return res
}

//go:generate goversioninfo -icon=ico/icon.ico -manifest=main.exe.manifest -arm=true
package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	// 插件示例

	// 以 Windows 为例，当用户在 主程序(StaticDeployment) 脚本中调用 `$命令名` 时，主程序会执行 `StaticDeployment.exe` 相同文件夹下的 `StaticDeployment_本插件编译后的文件名.exe` （如果主程序被重命名，则以新名称为准）。

	// 例如，用户在主程序脚本中调用 `$PluginDemo` 时，主程序会执行 `StaticDeployment_PluginDemo.exe` ，即本插件最终输出的执行文件名，并放在和主程序执行文件同一个文件夹下。

	// 插件程序在运行步骤的位置为下面的 `2.` 和 `3.`：
	// 1. 主程序备份正在修改的文件
	// 2. **主程序运行本插件程序**
	// 3. **本插件程序通过 fmt 将结果输出给主程序 或者 直接改好主程序通过参数提供的文件路径**
	// 4. 主程序根据插件程序改好的文件，或者插件程序提供的内容，进行编译等操作
	// 5. 主程序恢复备份的文件
	// 因此对用户的源文件进行覆盖修改是安全的也是通常做法。

	// 用户可以在主程序脚本中传递参数，例如 `$PluginDemo 参数1 参数2 参数...`，本插件会接收到 `参数0(本插件名) 参数1 参数2 参数...`。

	// - 如果需要输出日志信息，使用 `log` 包，内容会直接进行输出；
	// - 如果需要输出返回值，使用 `fmt` 包，内容会被主程序收取（如果用户在脚本中要求）。
	log.Println(os.Args[0], "测试插件开始运行……")

	// 当本插件需要问题时的处理方式
	if len(os.Args) <= 1 {
		log.Println(os.Args[0], "错误: 未指定参数。")
		// 如果本插件执行失败，退出时需返回 > 0 的值。主程序遇到非 0 的返回值会停止执行脚本。
		os.Exit(1)
		return
	}

	// 从主程序接收一个值
	var key = os.Args[1]

	// 假设这里有一些数据，根据用户的请求拿出来
	var data = map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	// 假设这里有个 HTML 代码片段，主程序需要拿这个 HTML 代码片段去合成整个 HTML 。
	var html = "<div>这是一个插件示例，根据 %KEY% ，我们找到了: %DATA%</div>"

	// 本插件将取到的数据填充到 HTML 代码片段中
	html = strings.ReplaceAll(html, "%KEY%", key)
	if value, ok := data[key]; ok {
		html = strings.ReplaceAll(html, "%DATA%", value)
	} else {
		html = strings.ReplaceAll(html, "%DATA%", "什么都没找到！")
	}

	// 将 HTML 代码片段返回给主程序。主程序会如实接收 fmt.Println 的内容，进行进一步操作。
	// 如果主程序脚本中没有设置接收数据，则 fmt.Println 会直接被输出（例如本插件是直接操作文件的程序）。
	fmt.Println(html)

	// 也可以由用户自行提供一个文件路径来由本插件打开或存储数据：
	// 如果插件的功能需要读取一个文件，可以将文件路径作为参数传递给本插件，本插件会读取文件内容并返回给主程序。
	if len(os.Args) > 2 {
		// 根据用户提供的路径读取文件
		var path = os.Args[2]
		data, err := os.ReadFile(path)
		if err != nil {
			log.Println(os.Args[0], "错误: 无法读取文件", path, err)
			os.Exit(1)
			return
		}
		// 将读取到的数据填充到 HTML 代码片段中
		html = strings.ReplaceAll(html, "%DATA%", string(data))
	}

	// 如果插件的功能需要将数据写入到一个文件，可以将文件路径作为参数传递给本插件，本插件会将数据写入到文件中。
	if len(os.Args) > 3 {
		var path = os.Args[3]
		// 例如，将 HTML 代码片段写入到文件中
		err := os.WriteFile(path, []byte(html), 0644)
		if err != nil {
			log.Println(os.Args[0], "错误: 无法写入文件", path, err)
			os.Exit(1)
			return
		}
	} else {
		// 可以判断用户是否给与另存为参数，如果没有，再通过输出来返回数据。
		fmt.Println(html)
	}

	// 本插件执行完毕，退出时需返回 0 的值。告诉主程序可以继续运行。
	os.Exit(0)

	// 在上面的程序中，用户在脚本中提供 `$PluginDemo key1` ，可以得到 `<div>这是一个插件示例，根据 key1 ，我们找到了: value1</div>` 的结果，并存储在一个自定义变量中。

	// 在上面的程序中，用户在脚本中提供 `$PluginDemo key2 1.html 1.html` ，并且 `1.html` 的内容为 `<div>这是一个插件示例，根据 %KEY% ，我们找到了: %DATA%</div>` , 可以得到 `<div>这是一个插件示例，根据 key2 ，我们找到了: value2</div>` 的结果，并且会将这个结果覆盖写入到 1.html 文件中。
}

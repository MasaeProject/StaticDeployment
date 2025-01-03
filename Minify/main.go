//go:generate goversioninfo -icon=ico/icon.ico -manifest=main.exe.manifest -arm=true
package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/json"
	"github.com/tdewolff/minify/v2/svg"
	"github.com/tdewolff/minify/v2/xml"
)

var mapMediaType map[string]string = map[string]string{
	"html": "text/html",
	"css":  "text/css",
	"js":   "application/javascript",
	"json": "application/json",
	"svg":  "image/svg+xml",
	"xml":  "application/xml",
}

var title string = "[代码压缩器] "

func StaticDeployment_Minify(cmd []string) ([2]int, error) {
	var cmdLen int = len(cmd)
	var path string = ""
	var outPath string = ""
	var mode string = ""
	var dataLen [2]int = [2]int{-1, -1}
	if cmdLen <= 1 {
		// path = srcPath
		return dataLen, fmt.Errorf(title + "请输入文件路径")
	} else if cmdLen >= 2 {
		path = cmd[1]
	}
	if cmdLen >= 3 {
		outPath = cmd[2]
	}
	if cmdLen >= 4 {
		mode = cmd[3]
	}
	// if len(outPath) == 0 {
	// 	outPath = path
	// }
	if len(mode) == 0 {
		var srcPathArr []string = strings.Split(path, ".")
		mode = srcPathArr[len(srcPathArr)-1]
	}
	if len(strings.Split(mode, "/")) == 1 {
		if _, ok := mapMediaType[mode]; ok {
			mode = mapMediaType[mode]
		}
	}
	sourceFileStat, err := os.Stat(path)
	if err != nil {
		return dataLen, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return dataLen, err
	}
	dataLen[0] = len(data)
	newCode, err := compressALL(string(data), mode)
	if err != nil {
		return dataLen, err
	}
	dataLen[1] = len(newCode)
	if len(outPath) == 0 {
		fmt.Println(string(newCode))
	} else {
		err = os.WriteFile(outPath, []byte(newCode), sourceFileStat.Mode())
		if err != nil {
			return dataLen, err
		}
	}
	return dataLen, nil
}

func compressALL(htmlContent string, mediatype string) (string, error) {
	m := minify.New()
	m.AddFunc(mapMediaType["css"], css.Minify)
	m.AddFunc(mapMediaType["html"], html.Minify)
	m.AddFunc(mapMediaType["svg"], svg.Minify)
	m.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
	m.AddFuncRegexp(regexp.MustCompile("[/+]json$"), json.Minify)
	m.AddFuncRegexp(regexp.MustCompile("[/+]xml$"), xml.Minify)

	return m.String(mediatype, htmlContent)
}

func main() {
	dataLen, err := StaticDeployment_Minify(os.Args)
	log.Printf("%s%d -> %d (E:%v)", title, dataLen[0], dataLen[1], err)
	if err != nil {
		os.Exit(1)
	}
}

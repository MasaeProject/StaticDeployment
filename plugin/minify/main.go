package codecompress

import (
	"fmt"
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

func InitWithCmd(cmd []string, srcPath string) ([2]int, error) {
	var cmdLen int = len(cmd)
	var path string = ""
	var outPath string = ""
	var mode string = ""
	var DataLen [2]int = [2]int{-1, -1}
	if cmdLen <= 1 {
		path = srcPath
	} else if cmdLen >= 2 {
		path = cmd[1]
	}
	if cmdLen >= 3 {
		mode = cmd[2]
	}
	if cmdLen >= 4 {
		outPath = cmd[3]
	}
	if len(outPath) == 0 {
		outPath = path
	}
	if len(mode) == 0 {
		var srcPathArr []string = strings.Split(path, ".")
		mode = srcPathArr[len(srcPathArr)-1]
	}
	if len(strings.Split(mode, "/")) == 1 {
		if _, ok := mapMediaType[mode]; ok {
			mode = mapMediaType[mode]
		}
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return DataLen, err
	}
	DataLen[0] = len(data)
	newCode, err := compressALL(string(data), mode)
	if err != nil {
		return DataLen, err
	}
	DataLen[1] = len(newCode)
	err = os.WriteFile(outPath, []byte(newCode), 0644)
	if err != nil {
		return DataLen, err
	}
	return DataLen, nil
}

func compressALL(htmlContent string, mediatype string) (string, error) {
	fmt.Println("Compressing", mediatype, "...")
	m := minify.New()
	m.AddFunc(mapMediaType["css"], css.Minify)
	m.AddFunc(mapMediaType["html"], html.Minify)
	m.AddFunc(mapMediaType["svg"], svg.Minify)
	m.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
	m.AddFuncRegexp(regexp.MustCompile("[/+]json$"), json.Minify)
	m.AddFuncRegexp(regexp.MustCompile("[/+]xml$"), xml.Minify)

	return m.String(mediatype, htmlContent)
}

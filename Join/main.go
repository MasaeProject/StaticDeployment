//go:generate goversioninfo -icon=ico/icon.ico -manifest=main.exe.manifest -arm=true
package main

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var title string = "[文件合并器] "

func StaticDeployment_Join(cmd []string) ([]int64, error) {
	var cmdLen int = len(cmd)
	var dataLen []int64 = make([]int64, cmdLen)
	if cmdLen <= 2 {
		return dataLen, fmt.Errorf(title + "需要输入至少两个文件路径")
	}
	var destPath string = os.Args[1]

	sourceFileStat, err := os.Stat(os.Args[2])
	if err != nil {
		return dataLen, err
	}
	var fileMode fs.FileMode = sourceFileStat.Mode()

	if err := os.MkdirAll(filepath.Dir(destPath), fileMode); err != nil {
		return dataLen, err
	}

	destinationFile, err := os.OpenFile(destPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, fileMode)
	if err != nil {
		return dataLen, err
	}
	defer destinationFile.Close()

	if destinationFileStat, err := destinationFile.Stat(); err == nil {
		dataLen[0] = destinationFileStat.Size()
	}

	var sourcePath []string = os.Args[2:]

	for i, sourcePath := range sourcePath {
		sourceFile, err := os.Open(sourcePath)
		if err != nil {
			return dataLen, err
		}
		defer sourceFile.Close()
		destinationFileStat, err := sourceFile.Stat()
		if err != nil {
			return dataLen, err
		}
		dataLen[i+1] = destinationFileStat.Size()
		_, err = io.Copy(destinationFile, sourceFile)
		if err != nil {
			return dataLen, err
		}
	}

	if destinationFileStat, err := destinationFile.Stat(); err == nil {
		dataLen[cmdLen-1] = destinationFileStat.Size()
	}

	return dataLen, nil
}

func main() {
	dataLen, err := StaticDeployment_Join(os.Args)
	var dataLenLen int = len(dataLen)
	var dataLenStrArr []string = make([]string, dataLenLen)
	for i, num := range dataLen {
		dataLenStrArr[i] = strconv.FormatInt(num, 10)
	}
	var total string = dataLenStrArr[dataLenLen-1]
	dataLenStrArr[0] = strings.Join(dataLenStrArr[:dataLenLen-1], " + ")
	log.Printf("%s%s = %s B  (E:%v)", title, dataLenStrArr[0], total, err)
	if err != nil {
		os.Exit(1)
		return
	}
}

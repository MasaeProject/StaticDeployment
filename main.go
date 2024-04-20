package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v2"
)

const (
	backupExtension = "StaticDeploymentBackup"
)

var (
	osName string
)

func main() {

	var solutions []Solution
	// var projects []Project

	osName = runtime.GOOS
	log.Println("StaticDeployment v0.0.1 for", osName)

	if len(os.Args) <= 1 {
		log.Println("必须指定一个配置文件路径。")
		return
	}
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Println("错误: 打开配置文件", os.Args[1], "失败:", err)
		return
	}
	defer file.Close()

	// buf := make([]byte, 512)
	// _, err = file.Read(buf)
	// if err != nil {
	// 	log.Println("错误: 读取配置文件失败:", err)
	// 	return
	// }

	var fileType byte = fileType(file)
	if fileType == 'j' {
		var decoder *json.Decoder = json.NewDecoder(file)
		err = decoder.Decode(&solutions)
		if err != nil {
			log.Println("错误: JSON 配置文件解析失败:", err)
			return
		}
	} else if fileType == 'y' {
		content, err := io.ReadAll(file)
		if err != nil {
			log.Println("错误: YAML 配置文件读取失败:", err)
			return
		}
		yaml.Unmarshal(content, &solutions)
	} else {
		log.Printf("错误: 未知的配置文件类型。")
		os.Exit(1)
		return
	}

	for i, solution := range solutions {
		log.Printf("开始处理: 解决方案 %d : %s\n", i+1, solution.Name)
		if runSolution(solution) {
			log.Printf("解决方案 %d : %s 处理完毕。\n", i+1, solution.Name)
		} else {
			log.Printf("解决方案 %d : %s 处理失败！\n", i+1, solution.Name)
		}
	}
}

func fileType(file *os.File) byte {
	var mime map[byte][]string = map[byte][]string{
		'j': {"json"},
		'y': {"yaml", "yml"},
	}
	filename := file.Name()
	extension := filepath.Ext(filename)
	for key, types := range mime {
		for _, vtype := range types {
			if extension == "."+vtype {
				return key
			}
		}
	}
	return 0
}

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

var (
	osName  string
	noChLog []BackupItem = []BackupItem{}
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

	var solutionLen = len(solutions)
	for i, solution := range solutions {
		log.Printf("开始处理: 解决方案 %d / %d : %s\n", i+1, solutionLen, solution.Name)
		if runSolution(solution) {
			log.Printf("解决方案 %d / %d : %s 处理完毕。\n", i+1, solutionLen, solution.Name)
		} else {
			log.Printf("解决方案 %d / %d : %s 处理失败！\n", i+1, solutionLen, solution.Name)
		}
	}

	if len(backupCache) > 0 {
		log.Println("警告: 不是所有备份文件都得到还原，请检查配置文件中备份和还原命令是否成对出现！未还原文件：")
		for i, bakInfo := range backupCache {
			var jobName Names = bakInfo.JobName
			log.Printf("%d  解决方案: %s  项目: %s  作业: %s  文件: %s\n", i+1, jobName.Solution, jobName.Project, jobName.Replace, bakInfo.SourceFile)
		}
	}
	if len(noChLog) > 0 {
		log.Println("警告: 不是所有文件都经过替换，请检查配置文件中的替换设置！替换前和替换后一样的文件有：")
		for i, bakInfo := range backupCache {
			var jobName Names = bakInfo.JobName
			log.Printf("%d  解决方案: %s  项目: %s  作业: %s  文件: %s\n", i+1, jobName.Solution, jobName.Project, jobName.Replace, bakInfo.SourceFile)
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

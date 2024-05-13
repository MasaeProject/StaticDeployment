//go:generate goversioninfo -icon=ico/icon.ico -manifest=main.exe.manifest -arm=true
package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

var (
	osName          string
	osExecFile      [4]string         = [4]string{} // 0: dir, 1: name, 2: ext, 3: realDir
	noChLog         []BackupItem      = []BackupItem{}
	totalIO         uint              = 0
	totalCMD        uint              = 0
	totalEXE        uint              = 0
	totalReplace    uint              = 0
	customVariables map[string]string = map[string]string{}
)

const backupExtension = "StaticDeploymentBackup"

func main() {
	var startTime time.Time = time.Now()
	var solutions []Solution
	// var projects []Project

	osName = runtime.GOOS
	log.Println("StaticDeployment v1.0.0 for", osName)
	log.Println("https://github.com/MasaeProject/StaticDeployment")

	getPaths()

	if len(os.Args) >= 2 && os.Args[1] == "-r" {
		rollback()
		return
	}

	var configPath = osExecFile[1] + ".yaml"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	} else if !Exists(configPath) {
		log.Println("错误: 必须指定一个配置文件。")
		os.Exit(1)
		return
	}
	if !Exists(configPath) {
		log.Println("错误: 配置文件路径不正确", configPath)
		os.Exit(1)
		return
	}
	file, err := os.Open(configPath)
	if err != nil {
		log.Println("错误: 打开配置文件", configPath, "失败:", err)
		os.Exit(1)
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
			os.Exit(1)
			return
		}
	} else if fileType == 'y' {
		content, err := io.ReadAll(file)
		if err != nil {
			log.Println("错误: YAML 配置文件读取失败:", err)
			os.Exit(1)
			return
		}
		yaml.Unmarshal(content, &solutions)
	} else {
		log.Printf("错误: 未知的配置文件类型。")
		os.Exit(1)
		return
	}
	var solutionLen = len(solutions)
	if solutionLen == 0 {
		log.Printf("错误: 配置文件格式不正确。")
		os.Exit(1)
		return
	}
	for i, solution := range solutions {
		log.Printf("开始处理: 解决方案 %d / %d : %s\n", i+1, solutionLen, solution.Name)
		var sTime time.Time = time.Now()
		if runSolution(solution) {
			log.Printf("解决方案 %d / %d : %s 处理完毕，用时 %.2f 秒。\n", i+1, solutionLen, solution.Name, time.Since(sTime).Seconds())
		} else {
			log.Printf("解决方案 %d / %d : %s 处理失败！用时 %.2f 秒。\n", i+1, solutionLen, solution.Name, time.Since(sTime).Seconds())
			return
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
		log.Println("警告: 替换前和替换后内容一样的文件有：")
		for i, bakInfo := range noChLog {
			var jobName Names = bakInfo.JobName
			log.Printf("%d  解决方案: %s  项目: %s  作业: %s  文件: %s\n", i+1, jobName.Solution, jobName.Project, jobName.Replace, bakInfo.SourceFile)
		}
	}
	duration := time.Since(startTime)
	log.Printf("处理成功，用时 %.2f 秒。替换项: %d ; 文件操作: %d ; 执行命令: %d (外部命令: %d )\n", duration.Seconds(), totalReplace, totalIO, totalCMD, totalEXE)
}

func getPaths() {
	var osFileNameArr []string = strings.Split(os.Args[0], string(filepath.Separator))
	osExecFile[0] = strings.Join(osFileNameArr[:len(osFileNameArr)-1], string(filepath.Separator))
	if len(osExecFile[0]) == 0 {
		osExecFile[0] = "."
	}
	osExecFile[1] = osFileNameArr[len(osFileNameArr)-1]
	osFileNameArr = strings.Split(osExecFile[1], ".")
	if len(osFileNameArr) > 1 {
		osExecFile[1] = osFileNameArr[0]
		osExecFile[2] = osFileNameArr[1]
		if len(osExecFile[2]) > 0 {
			osExecFile[2] = "." + osExecFile[2]
		}
	}
	execPath, err := os.Executable()
	if err == nil {
		osExecFile[3] = execPath
		realPath, err := filepath.EvalSymlinks(execPath)
		if err == nil {
			osExecFile[3] = realPath
		}
	}
	osFileNameArr = strings.Split(osExecFile[3], string(filepath.Separator))
	osExecFile[3] = strings.Join(osFileNameArr[:len(osFileNameArr)-1], string(filepath.Separator))
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

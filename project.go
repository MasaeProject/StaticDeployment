package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

func runProject(project Project) bool {
	log.Println("开始处理项目:", project.Name)
	if project.Source == "" {
		log.Println("错误：你必须为项目指定一个 source 。")
		return false
	}
	loadFileText, err := os.ReadFile(project.Source)
	var f FileData = FileData{Path: project.Source, LoadString: string(loadFileText), LoadSize: len(loadFileText)}
	if err != nil {
		log.Printf("错误：打开文件 %s 失败：%s\n", project.Source, err)
		return false
	} else {
		log.Printf("已加载文件 %s (%d B)\n", project.Source, len(loadFileText))
	}
	if runReplace(project.Replace, f) {
		runExec(project.Exec, project.Source, "")
	}
	return true
}

func ensureDir(filePath string) bool {
	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Printf("错误：创建文件夹 %s 失败：%s\n", dir, err)
			return false
		}
	}
	return true
}

func runReplace(replace []ReplaceItem, f FileData) bool {
	if len(replace) == 0 {
		return true
	}
	var isOK = true
	for _, item := range replace {
		f = runReplaceDetail(item.Replace, f)
		if ensureDir(item.To) {
			err := os.WriteFile(item.To, []byte(f.NewString), 0666)
			if err != nil {
				isOK = false
				log.Printf("错误：写入文件 %s 失败：%s (%d B -> %d B)\n", item.To, err, f.LoadSize, f.NewSize)
			} else {
				log.Printf("已写入文件 %s (%d B -> %d B)\n", item.To, f.LoadSize, f.NewSize)
				runExec(item.Exec, f.Path, item.To)
			}
		}
	}
	return isOK
}

func runExec(exec [][]string, srcPath string, toPath string) {
	if len(exec) == 0 {
		return
	}
	var isRun bool = false
	var defaultCmd []string = []string{}
	for _, nowExec := range exec {
		var useOS string = nowExec[0]
		var cmd []string = nowExec[1:]
		for i, c := range cmd {
			if c == "$source$" || c == "$src$" {
				cmd[i] = srcPath
			} else if c == "$to$" {
				cmd[i] = toPath
			}
		}
		if len(useOS) == 0 {
			defaultCmd = cmd
		} else if useOS == osName {
			isRun = true
			runCMD(cmd)
		}
	}
	if !isRun && defaultCmd != nil {
		runCMD(defaultCmd)
	}
}

func runCMD(cmd []string) {
	log.Println("运行命令:", strings.Join(cmd, " "))
	ex := exec.Command(cmd[0], cmd[1:]...)
	ex.Stdout = os.Stdout
	ex.Stderr = os.Stderr
	err := ex.Run()
	if exitError, ok := err.(*exec.ExitError); ok {
		if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
			fmt.Printf("命令退出代码: %d\n", status.ExitStatus())
		}
	} else if err != nil {
		log.Println("错误：执行命令失败：", err)
	} else {
		log.Println("命令运行成功。")
	}
}

func runReplaceDetail(replace []ReplaceDetail, f FileData) FileData {
	f.NewString = f.LoadString
	if len(replace) > 0 {
		for _, detail := range replace {
			if detail.Num < 0 {
				f.NewString = strings.ReplaceAll(f.LoadString, detail.Old, detail.New)
			} else {
				f.NewString = strings.Replace(f.LoadString, detail.Old, detail.New, detail.Num)
			}
		}
	}
	f.NewSize = len([]byte(f.NewString))
	return f
}

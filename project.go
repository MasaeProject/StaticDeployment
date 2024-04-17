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
	if project.PreRun != nil {
		runExec(*project.PreRun, project.Source, "")
	}
	if runReplace(project.Replace, f) && project.Run != nil {
		runExec(*project.Run, project.Source, "")
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
			if item.PreRun != nil {
				runExec(*item.PreRun, f.Path, item.To)
			}
			err := os.WriteFile(item.To, []byte(f.NewString), 0666)
			if err != nil {
				isOK = false
				log.Printf("错误：写入文件 %s 失败：%s (%d B -> %d B)\n", item.To, err, f.LoadSize, f.NewSize)
			} else {
				log.Printf("已写入文件 %s (%d B -> %d B)\n", item.To, f.LoadSize, f.NewSize)
				if item.Run != nil {
					runExec(*item.Run, f.Path, item.To)
				}
			}
		}
	}
	return isOK
}

func runExec(run Run, srcPath string, toPath string) {
	var cmds [][]string = [][]string{}
	if osName == "windows" && run.Windows != nil && len(*run.Windows) > 0 {
		cmds = *run.Windows
	} else if osName == "linux" && run.Linux != nil && len(*run.Linux) > 0 {
		cmds = *run.Linux
	} else if osName == "darwin" && run.Darwin != nil && len(*run.Darwin) > 0 {
		cmds = *run.Darwin
	} else if run.Default != nil && len(*run.Default) > 0 {
		cmds = *run.Default
	} else {
		return
	}
	for _, cmd := range cmds {
		var noEmbCmd = false
		if len(cmd) == 3 {
			var err error = nil
			switch cmd[0] {
			case "$CP":
				err = Copy(cmd[1], cmd[2])
			case "$MV":
				err = Move(cmd[1], cmd[2])
			case "$SMV":
				err = MoveSecure(cmd[1], cmd[2])
			case "$RM":
				err = Remove(cmd[1])
			case "$SRM":
				err = RemoveSecure(cmd[1])
			case "$REN":
				err = Rename(cmd[1], cmd[2])
			default:
				noEmbCmd = true
			}
			if err != nil {
				log.Printf("错误: 文件操作 %s \"%s\" \"%s\" 失败: %s\n", cmd[0], cmd[1], cmd[2], err)
			}
		} else {
			noEmbCmd = true
		}
		if noEmbCmd {
			for i, c := range cmd {
				if c == "$SOURCE" || c == "$SRC" {
					cmd[i] = srcPath
				} else if c == "$TO" {
					cmd[i] = toPath
				}
			}
			runCMD(cmd)
		}
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

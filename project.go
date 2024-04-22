package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

func runSolution(solution Solution) bool {
	var names Names = Names{Solution: solution.Name, Project: "", Replace: ""}
	if solution.PreRun != nil {
		log.Println("运行解决方案", solution.Name, "的预处理命令:")
		if !runExec(*solution.PreRun, "", names) {
			return false
		}
	}
	var projectsLen int = len(solution.Projects)
	for i, project := range solution.Projects {
		log.Printf("开始处理: 解决方案: %s  项目 %d / %d : %s\n", names.Solution, i+1, projectsLen, project.Name)
		var sTime time.Time = time.Now()
		if runProject(project, names) {
			log.Printf("项目 %d / %d : %s 处理完毕，用时 %.2f 秒。\n", i+1, projectsLen, project.Name, time.Since(sTime).Seconds())
		} else {
			log.Printf("项目 %d / %d : %s 处理失败！用时 %.2f 秒。\n", i+1, projectsLen, project.Name, time.Since(sTime).Seconds())
			return false
		}
	}
	if solution.Run != nil {
		log.Println("运行解决方案", solution.Name, "的后处理命令:")
		if !runExec(*solution.PreRun, "", names) {
			return false
		}
	}
	return true
}

func runProject(project Project, names Names) bool {
	names.Project = project.Name
	if project.Source == "" {
		log.Println("错误：你必须为项目指定一个 source 。")
		return false
	}
	project.Source = CleanPath(project.Source)
	totalIO++
	loadFileText, err := os.ReadFile(project.Source)
	var f FileData = FileData{Path: project.Source, LoadString: string(loadFileText), LoadSize: len(loadFileText)}
	if err != nil {
		log.Printf("错误：打开文件 %s 失败：%s\n", project.Source, err)
		return false
	} else {
		log.Printf("已加载文件 %s (%d B)\n", project.Source, len(loadFileText))
	}
	if project.PreRun != nil {
		log.Println("运行项目", project.Name, "的预处理命令:")
		if !runExec(*project.PreRun, project.Source, names) {
			return false
		}
	}
	var replaceLen int = len(project.Replace)
	for i, item := range project.Replace {
		log.Printf("开始处理: 解决方案: %s  项目: %s  作业 %d / %d : %s\n", names.Solution, project.Name, i+1, replaceLen, item.Name)
		var sTime time.Time = time.Now()
		if runJob(item, f, names) {
			log.Printf("作业 %d / %d : %s 处理完毕，用时 %.2f 秒。\n", i+1, replaceLen, item.Name, time.Since(sTime).Seconds())
		} else {
			log.Printf("作业 %d / %d : %s 处理失败！用时 %.2f 秒。\n", i+1, replaceLen, item.Name, time.Since(sTime).Seconds())
			return false
		}
	}
	if project.Run != nil {
		log.Println("运行项目", project.Name, "的后处理命令:")
		if !runExec(*project.Run, project.Source, names) {
			return false
		}
	}
	return true
}

func runJob(item ReplaceItem, f FileData, names Names) bool {
	names.Replace = item.Name
	if len(item.Replace) == 0 {
		log.Println("警告: 没有指定替换方案")
	}
	f = runReplaceDetail(item.Replace, f)
	// var bak string = f.Path + "." + backupExtension
	if item.PreRun != nil {
		log.Println("运行作业", item.Name, "的预处理命令:")
		if !runExec(*item.PreRun, f.Path, names) {
			return false
		}
	}
	// var err error = copyFile(f.Path, bak)
	// if err != nil {
	// 	log.Printf("错误: 备份文件 %s 到 %s 失败: %s\n", f.Path, bak, err)
	// 	return false
	// }
	totalIO++
	if err := os.WriteFile(f.Path, []byte(f.NewString), 0666); err != nil {
		log.Printf("错误: 写入文件 %s 失败：%s (%d B -> %d B)\n", f.Path, err, f.LoadSize, f.NewSize)
		return false
	} else {
		var noCh = ""
		if f.LoadString == f.NewString {
			var noChItem BackupItem = BackupItem{SourceFile: f.Path, JobName: names}
			noChLog = append(noChLog, noChItem)
			noCh = " (无变化)"
			// fmt.Println("========== 无变化 ==========")
			// fmt.Printf("文件: %v\n", noChItem)
			// fmt.Printf("方案: %v\n", item.Replace)
			// fmt.Println(f.LoadString)
			// fmt.Println("====================")
		}
		log.Printf("已写入文件 %s (%d B -> %d B)%s\n", f.Path, f.LoadSize, f.NewSize, noCh)
		if item.Run != nil {
			log.Println("运行作业", item.Name, "的后处理命令:")
			if !runExec(*item.Run, f.Path, names) {
				return false
			}
		}
	}
	// err = removeFile(f.Path)
	// if err != nil {
	// 	log.Printf("警告: 删除临时文件 %s 失败: %s\n", f.Path, err)
	// }
	// err = RenamePath(bak, f.Path)
	// if err != nil {
	// 	log.Printf("错误: 恢复文件 %s 到 %s 失败: %s\n", bak, f.Path, err)
	// 	return false
	// }
	return true
}

func runExec(run Run, srcPath string, names Names) bool {
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
		log.Println("错误: 未找到适用于当前操作系统的命令。")
		return false
	}
	srcPath = CleanPath(srcPath)
	var cmdsLen int = len(cmds)
	totalCMD++
	var dir string = ""
	for cmdsI, cmd := range cmds {
		var noEmbCmd = false
		var cmdLen int = len(cmd)
		for cmdI, c := range cmd {
			var nKey string = "$SRC"
			c = CleanPath(c)
			var pathArr []string = strings.Split(srcPath, string(filepath.Separator))
			var pathArrLen = len(pathArr)
			var fileFullName string = pathArr[pathArrLen-1]
			var dirPath string = ""
			if pathArrLen > 1 {
				fileFullName = pathArr[pathArrLen-1]
				pathArr = pathArr[:pathArrLen-1]
				dirPath = strings.Join(pathArr, string(filepath.Separator))
			}
			var fileNameArr []string = strings.Split(fileFullName, ".")
			var fileNameArrLen = len(fileNameArr)
			var extName string = fileNameArr[fileNameArrLen-1]
			var fileName string = fileNameArr[0]
			if fileNameArrLen > 2 {
				fileNameArr = fileNameArr[:fileNameArrLen-1]
				fileName = strings.Join(fileNameArr, ".")
			} else if fileNameArrLen == 1 {
				extName = ""
			}
			nKey = "$SOLUTION"
			if strings.Contains(c, nKey) {
				c = strings.ReplaceAll(c, nKey, names.Solution)
			}
			nKey = "$PROJECT"
			if strings.Contains(c, nKey) {
				c = strings.ReplaceAll(c, nKey, names.Project)
			}
			nKey = "$JOBNAME"
			if strings.Contains(c, nKey) {
				c = strings.ReplaceAll(c, nKey, names.Replace)
			}
			nKey = "$SRCFILE"
			if strings.Contains(c, nKey) {
				c = strings.ReplaceAll(c, nKey, fileFullName)
			}
			nKey = "$SRCNAME"
			if strings.Contains(c, nKey) {
				c = strings.ReplaceAll(c, nKey, fileName)
			}
			nKey = "$SRCEXT"
			if strings.Contains(c, nKey) {
				if IsDirectory(srcPath) == 0 {
					c = strings.ReplaceAll(c, nKey, extName)
				} else {
					c = strings.ReplaceAll(c, nKey, "")
				}
			}
			nKey = "$SRCDIRNAME"
			if strings.Contains(c, nKey) {
				c = strings.ReplaceAll(c, nKey, pathArr[len(pathArr)-1])
			}
			nKey = "$SRCDIR"
			if strings.Contains(c, nKey) {
				c = strings.ReplaceAll(c, nKey, dirPath)
			}
			nKey = "$SRC"
			if strings.Contains(c, nKey) {
				c = strings.ReplaceAll(c, nKey, srcPath)
			}
			// fmt.Println("SRCFILE", fileFullName, "SRCNAME", fileName, "SRCEXT", extName, "SRCDIRNAME", pathArr[len(pathArr)-1], "SRCDIR", dirPath, "SRC", srcPath)
			cmd[cmdI] = CleanPath(c)
		}

		if len(cmd) == 0 {
			continue
		}
		log.Printf("运行命令 %d / %d : %s\n", cmdsI+1, cmdsLen, strings.Join(cmd, " "))
		var err error = nil
		var isOK bool = true
		switch cmd[0] {
		case "$CMDDIR":
			if cmdLen == 1 {
				dir = ""
			} else if cmdLen >= 2 {
				dir = cmd[1]
			}
		case "$BAK":
			if cmdLen == 1 || (cmdLen == 2 && len(cmd[1]) == 0) {
				isOK = backup(srcPath, names)
			}
			if cmdLen >= 2 {
				isOK = backup(cmd[1], names)
			}
		case "$RES":
			var resCmd Names = names
			if cmdLen >= 2 && len(cmd[1]) > 0 {
				resCmd.Solution = cmd[1]
			}
			if cmdLen >= 3 && len(cmd[2]) > 0 {
				resCmd.Project = cmd[2]
			}
			if cmdLen >= 3 && len(cmd[3]) > 0 {
				resCmd.Replace = cmd[3]
			}
			if cmdLen == 2 {
				isOK = restoreSolution(resCmd.Solution)
			} else if cmdLen >= 3 {
				isOK = restoreProject(resCmd.Solution, resCmd.Project)
			} else if cmdLen == 1 || cmdLen == 4 {
				isOK = restoreJob(resCmd.Solution, resCmd.Project, resCmd.Replace)
			}
		case "$CP":
			if cmdLen >= 3 {
				err = Copy(cmd[1], cmd[2])
			} else if cmdLen == 2 {
				err = Copy(srcPath, cmd[1])
			}
		case "$MV":
			if cmdLen >= 3 {
				err = Move(cmd[1], cmd[2])
			} else if cmdLen == 2 {
				err = Move(srcPath, cmd[1])
			}
		case "$SMV":
			if cmdLen >= 3 {
				err = MoveSecure(cmd[1], cmd[2])
			} else if cmdLen == 2 {
				err = MoveSecure(srcPath, cmd[1])
			}
		case "$RM":
			err = Remove(cmd[1])
		case "$SRM":
			err = RemoveSecure(cmd[1])
		case "$REN":
			if cmdLen >= 3 {
				err = RenamePath(cmd[1], cmd[2])
			} else if cmdLen == 2 {
				err = RenamePath(srcPath, cmd[1])
			}
		default:
			noEmbCmd = true
		}
		if err != nil {
			log.Printf("错误: 文件操作 %s \"%s\" \"%s\" 失败: %s\n", cmd[0], cmd[1], cmd[2], err)
			return false
		} else if !isOK {
			return false
		}

		if noEmbCmd {
			if !runCMD(cmd, dir) {
				return false
			}
		}
	}
	return true
}

func runCMD(cmd []string, dir string) bool {
	totalEXE++
	ex := exec.Command(cmd[0], cmd[1:]...)
	if len(dir) > 0 {
		ex.Dir = dir
	}
	ex.Stdout = os.Stdout
	ex.Stderr = os.Stderr
	err := ex.Run()
	if exitError, ok := err.(*exec.ExitError); ok {
		if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
			fmt.Printf("错误: 命令退出代码: %d\n", status.ExitStatus())
			return false
		}
	} else if err != nil {
		log.Println("错误：执行命令失败：", err)
		return false
	} else {
		log.Println("命令运行成功。")
		return true
	}
	return true
}

func runReplaceDetail(replace []ReplaceDetail, f FileData) FileData {
	f.NewString = f.LoadString
	if len(replace) > 0 {
		for _, detail := range replace {
			if detail.Num < 0 {
				var count int = strings.Count(f.LoadString, detail.Old)
				f.NewString = strings.ReplaceAll(f.LoadString, detail.Old, detail.New)
				totalReplace = totalReplace + uint(count)
			} else {
				f.NewString = strings.Replace(f.LoadString, detail.Old, detail.New, detail.Num)
				totalReplace++
			}
		}
	}
	f.NewSize = len([]byte(f.NewString))
	return f
}

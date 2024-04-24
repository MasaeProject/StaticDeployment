package main

import (
	"log"
	"os"
	"strings"
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
		log.Printf("开始处理: 解决方案: %s  项目: %s  替换 %d / %d : %s\n", names.Solution, project.Name, i+1, replaceLen, item.Name)
		var sTime time.Time = time.Now()
		if runJob(item, f, names) {
			log.Printf("替换 %d / %d : %s 处理完毕，用时 %.2f 秒。\n", i+1, replaceLen, item.Name, time.Since(sTime).Seconds())
		} else {
			log.Printf("替换 %d / %d : %s 处理失败！用时 %.2f 秒。\n", i+1, replaceLen, item.Name, time.Since(sTime).Seconds())
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
		log.Println("运行替换", item.Name, "的预处理命令:")
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
	}
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
		log.Println("运行替换", item.Name, "的后处理命令:")
		if !runExec(*item.Run, f.Path, names) {
			return false
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

func runReplaceDetail(replace []ReplaceDetail, f FileData) FileData {
	var newString string = f.LoadString
	if len(replace) > 0 {
		for _, detail := range replace {
			var count int = strings.Count(newString, detail.Old)
			if detail.Num < 0 {
				newString = strings.ReplaceAll(newString, detail.Old, detail.New)
			} else {
				newString = strings.Replace(newString, detail.Old, detail.New, detail.Num)
				if count > detail.Num {
					count = detail.Num
				}
			}
			totalReplace = totalReplace + uint(count)
		}
	}
	f.NewString = newString
	f.NewSize = len([]byte(newString))
	return f
}

package main

import (
	"fmt"
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
	if solution.Projects != nil {
		var projects []Project = *solution.Projects
		var projectsLen int = len(projects)
		for i, project := range projects {
			log.Printf("开始处理: 解决方案: %s  项目 %d / %d : %s\n", names.Solution, i+1, projectsLen, project.Name)
			var sTime time.Time = time.Now()
			if runProject(project, names) {
				log.Printf("项目 %d / %d : %s 处理完毕，用时 %.2f 秒。\n", i+1, projectsLen, project.Name, time.Since(sTime).Seconds())
			} else {
				log.Printf("项目 %d / %d : %s 处理失败！用时 %.2f 秒。\n", i+1, projectsLen, project.Name, time.Since(sTime).Seconds())
				return false
			}
		}
	}
	if solution.Run != nil {
		log.Println("运行解决方案", solution.Name, "的后处理命令:")
		if !runExec(*solution.Run, "", names) {
			return false
		}
	}
	return true
}

func runProject(project Project, names Names) bool {
	names.Project = project.Name
	var source string = ""
	var f FileData = FileData{}
	if project.Source != nil {
		source = *project.Source
		source = CleanPath(source)
		totalIO++
		loadFileText, err := os.ReadFile(source)
		f = FileData{Path: source, LoadString: string(loadFileText), LoadSize: len(loadFileText)}
		if err != nil {
			log.Printf("错误：打开文件 %s 失败：%s\n", source, err)
			return false
		} else {
			log.Printf("已加载文件 %s (%d B)\n", source, len(loadFileText))
		}
		project.Source = &source
	}

	if project.PreRun != nil {
		log.Println("运行项目", project.Name, "的预处理命令:")
		if !runExec(*project.PreRun, source, names) {
			return false
		}
	}
	if project.Replaces != nil {
		var replaces []ReplaceItem = *project.Replaces
		var replaceLen int = len(replaces)
		for i, item := range replaces {
			log.Printf("开始处理: 解决方案: %s  项目: %s  替换 %d / %d : %s\n", names.Solution, project.Name, i+1, replaceLen, item.Name)
			var sTime time.Time = time.Now()
			if runJob(item, f, names) {
				log.Printf("替换任务 %d / %d : %s 处理完毕，用时 %.2f 秒。\n", i+1, replaceLen, item.Name, time.Since(sTime).Seconds())
			} else {
				log.Printf("替换任务 %d / %d : %s 处理失败！用时 %.2f 秒。\n", i+1, replaceLen, item.Name, time.Since(sTime).Seconds())
				return false
			}
		}
	}
	if project.Run != nil {
		log.Println("运行项目", project.Name, "的后处理命令:")
		if !runExec(*project.Run, source, names) {
			return false
		}
	}
	return true
}

func runJob(item ReplaceItem, f FileData, names Names) bool {
	names.Replace = item.Name
	// var bak string = f.Path + "." + backupExtension
	if item.PreRun != nil {
		log.Println("运行替换", item.Name, "的预处理命令:")
		if !runExec(*item.PreRun, f.Path, names) {
			return false
		}
	}
	if len(item.Items) == 0 {
		log.Println("警告: 没有指定替换方案")
	}
	f = runReplaceDetail(item.Items, f)
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
	var replaceLen int = len(replace)
	if replaceLen > 0 {
		for i, detail := range replace {
			var newStr = detail.New
			for key, val := range customVariables {
				newStr = strings.ReplaceAll(newStr, "$"+key, val)
			}
			if len(newStr) > 8 && newStr[:8] == "$IMPORT=" {
				var importFiles string = newStr[8:]
				newStr = ""
				var importFilesList []string = strings.Split(importFiles, ",")
				for _, importFile := range importFilesList {
					importFile = CleanPath(importFile)
					totalIO++
					loadFileText, err := os.ReadFile(importFile)
					if err != nil {
						log.Printf("错误：导入文件 %s 失败：%s\n", importFile, err)
						os.Exit(1)
					} else {
						log.Printf("已加载文件 %s (%d B)\n", importFile, len(loadFileText))
						newStr += string(loadFileText)
					}
				}
			}
			var count int = strings.Count(newString, detail.Old)
			var num = 1
			if detail.Num != nil {
				num = *detail.Num
			}
			if num < 0 {
				newString = strings.ReplaceAll(newString, detail.Old, newStr)
			} else {
				newString = strings.Replace(newString, detail.Old, newStr, num)
				if count > num {
					count = num
				}
			}
			log.Printf("替换项 %d / %d : %s -> %s (%d -> %d)", i+1, replaceLen, trimString(detail.Old), trimString(newStr), len(detail.Old), len(newStr))
			if detail.Old == newStr {
				log.Println("警告: 找不到替换项或替换前后内容一样")
			}
			totalReplace = totalReplace + uint(count)
			// log.Printf("newString =  %s \n", newString)
		}
	}
	f.NewString = newString
	f.NewSize = len([]byte(newString))
	return f
}

func trimString(input string) string {
	runes := []rune(strings.ReplaceAll(strings.ReplaceAll(input, "\r", ""), "\n", ""))
	length := len(runes)
	if length > 20 {
		return fmt.Sprintf("\"%s\" ... \"%s\"", string(runes[:10]), string(runes[length-10:]))
	}
	return "\"" + input + "\""
}

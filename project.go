package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func runProject(project Project) bool {
	log.Println("开始处理项目:", project.Name)
	if project.Source == "" {
		log.Println("错误：你必须为项目指定一个 source 。")
		return false
	}
	loadFileText, err := os.ReadFile(project.Source)
	var f FileData = FileData{LoadString: string(loadFileText), LoadSize: len(loadFileText)}
	if err != nil {
		log.Printf("错误：打开文件 %s 失败：%s\n", project.Source, err)
		return false
	} else {
		log.Printf("已加载文件 %s (%d B)\n", project.Source, len(loadFileText))
	}
	runReplace(project.Replace, f)
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

func runReplace(replace []ReplaceItem, f FileData) {
	if len(replace) == 0 {
		return
	}
	for _, item := range replace {
		f = runReplaceDetail(item.Replace, f)
		if ensureDir(item.To) {
			err := os.WriteFile(item.To, []byte(f.NewString), 0666)
			if err != nil {
				log.Printf("错误：写入文件 %s 失败：%s (%d B -> %d B)\n", item.To, err, f.LoadSize, f.NewSize)
			} else {
				log.Printf("已写入文件 %s (%d B -> %d B)\n", item.To, f.LoadSize, f.NewSize)
			}
		}
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

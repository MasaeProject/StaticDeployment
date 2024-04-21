package main

import (
	"log"
	"strings"
)

const backupExtension = "StaticDeploymentBackup"

var backupCache []BackupItem = []BackupItem{}

func isBackupPath(srcPath string) (bool, string) {
	var srcExtArr []string = strings.Split(srcPath, ".")
	if len(srcExtArr) > 0 && srcExtArr[len(srcExtArr)-1] == backupExtension {
		return true, strings.Join(srcExtArr[:len(srcExtArr)-1], ".")
	} else {
		return false, srcPath + "." + backupExtension
	}
}

func backup(srcPath string, names Names) bool {
	var isBak, bakPath = isBackupPath(srcPath)
	log.Printf("备份文件 %d : %s\n", len(backupCache)+1, bakPath)
	if isBak {
		log.Printf("错误: %s 已经是备份文件。\n", srcPath)
		return false
	}
	if !Exists(srcPath) {
		log.Printf("错误: 要备份的源文件 %s 不存在。\n", srcPath)
		return false
	}
	if Exists(bakPath) {
		if err := Remove(bakPath); err != nil {
			log.Printf("错误: 删除备份文件 %s 失败: %s\n", bakPath, err)
			return false
		}
	}
	if err := Copy(srcPath, bakPath); err != nil {
		log.Printf("错误: 从 %s 备份到 %s 失败: %s\n", srcPath, bakPath, err)
		return false
	}
	var cacheInfo BackupItem = BackupItem{SourceFile: srcPath, JobName: names}
	backupCache = append(backupCache, cacheInfo)
	return true
}

func restore(bakPath string, i int) bool {
	var isBak, srcPath = isBackupPath(bakPath)
	log.Printf("还原文件 %d / %d : %s\n", i, len(backupCache), srcPath)
	if !isBak {
		log.Printf("错误: %s 不是备份文件。\n", bakPath)
		return false
	}
	if !Exists(bakPath) {
		log.Printf("错误: 要还原的备份文件 %s 不存在。\n", bakPath)
		return false
	}
	if Exists(srcPath) {
		if err := Remove(srcPath); err != nil {
			log.Printf("错误: 删除源文件 %s 失败: %s\n", srcPath, err)
			return false
		}
	}
	if err := Copy(bakPath, srcPath); err != nil {
		log.Printf("错误: 从 %s 还原到 %s 失败: %s\n", bakPath, srcPath, err)
		return false
	}
	return true
}

func restoreSolution(solutionName string) bool {
	for i, item := range backupCache {
		if item.JobName.Solution == solutionName {
			var bakPath = item.SourceFile + "." + backupExtension
			rmbackupCache(i)
			return restore(bakPath, i)
		}
	}
	log.Printf("错误: 未找到解决方案 %s 的备份文件。\n", solutionName)
	return false
}

func restoreProject(solutionName, projectName string) bool {
	for i, item := range backupCache {
		if item.JobName.Solution == solutionName && item.JobName.Project == projectName {
			var bakPath = item.SourceFile + "." + backupExtension
			rmbackupCache(i)
			return restore(bakPath, i)
		}
	}
	log.Printf("错误: 未找到解决方案 %s 项目 %s 的备份文件。\n", solutionName, projectName)
	return false
}

func restoreJob(solutionName, projectName, jobName string) bool {
	for i, item := range backupCache {
		if item.JobName.Solution == solutionName && item.JobName.Project == projectName && item.JobName.Replace == jobName {
			var bakPath = item.SourceFile + "." + backupExtension
			rmbackupCache(i)
			return restore(bakPath, i)
		}
	}
	log.Printf("错误: 未找到解决方案 %s 项目 %s 作业 %s 的备份文件。\n", solutionName, projectName, jobName)
	return false
}

func rmbackupCache(index int) {
	if index > 0 && index <= len(backupCache) {
		backupCache = append(backupCache[:index-1], backupCache[index:]...)
	}
}
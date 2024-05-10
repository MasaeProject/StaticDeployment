package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func rollback() {
	var rootDir string = "."
	var totalDirs []uint = []uint{0, 0}
	if len(os.Args) >= 3 {
		rootDir = os.Args[2]
	}
	log.Println("回滚模式: 扫描文件夹:", rootDir)
	var bak string = "." + backupExtension
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		totalDirs[0]++
		if err != nil {
			return err
		}
		if strings.HasSuffix(info.Name(), bak) {
			newName := strings.TrimSuffix(path, bak)
			log.Printf("%s -> %s\n", path, newName)
			// 檢查不帶bak的檔案/資料夾是否存在
			if _, err := os.Stat(newName); err == nil {
				// 如果存在，則刪除
				err := os.Remove(newName)
				if err != nil {
					return fmt.Errorf("无法删除已存在的文件 %s: %w", newName, err)
				}
			}
			// 重新命名bak檔案/資料夾
			err = os.Rename(path, newName)
			if err != nil {
				return fmt.Errorf("无法重命名 %s to %s: %w", path, newName, err)
			}
			totalDirs[1]++
		}
		return nil
	})
	if err != nil {
		log.Println("扫描文件夹失败:", err)
	}
	log.Printf("回滚结束。已操作文件: %d / %d\n", totalDirs[1], totalDirs[0])
}

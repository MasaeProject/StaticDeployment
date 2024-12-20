package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// isHiddenFile 判斷給定的檔案或目錄是否為隱藏檔案或隱藏目錄
// 引數:
// - path: 檔案的完整路徑
// - info: 檔案的 os.FileInfo 資訊
// 返回值:
// - bool: 如果是隱藏檔案或目錄，返回 true，否則返回 false
func isHiddenFile(path string, info os.FileInfo) bool {
	if path == "." {
		return false // 保留當前目錄
	}
	// Linux 隱藏檔案以 '.' 開頭
	if strings.HasPrefix(info.Name(), ".") {
		return true
	}
	// Windows 隱藏檔案檢測，可以使用檔案屬性或簡單名稱檢查
	if info.Name() == "" || info.Name()[0] == '.' {
		return true
	}
	return false
}

// rollback 掃描指定目錄並回滾檔案（重新命名備份檔案）
// 掃描時跳過隱藏檔案/目錄和 node_modules 目錄。
// 引數:
// - 無直接引數，從命令列讀取根目錄引數。
func rollback() {
	var rootDir string = "."            // 預設掃描當前目錄
	var totalDirs []uint = []uint{0, 0} // totalDirs[0]: 總檔案數, totalDirs[1]: 已處理檔案數
	if len(os.Args) >= 3 {
		rootDir = os.Args[2]
	}
	log.Println("回滚模式: 扫描文件夹:", rootDir)
	var bak string = "." + backupExtension // 備份檔案的副檔名

	// 遍歷目錄和檔案
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err // 如果訪問檔案出錯，返回錯誤
		}
		// 跳過隱藏檔案和隱藏目錄
		if isHiddenFile(path, info) {
			// log.Printf("跳過隱藏檔案/目錄: %s\n", path)
			if info.IsDir() {
				return filepath.SkipDir // 跳過隱藏目錄
			}
			return nil
		}
		// 跳過 node_modules 目錄
		if info.IsDir() && info.Name() == "node_modules" {
			// log.Printf("跳過目錄: %s\n", path)
			return filepath.SkipDir // 跳過 node_modules 目錄
		}

		totalDirs[0]++ // 檔案數累加
		// 檢查是否為備份檔案，並回滾
		if strings.HasSuffix(info.Name(), bak) {
			newName := strings.TrimSuffix(path, bak) // 去掉備份副檔名
			log.Printf("%s -> %s\n", path, newName)
			// 如果目標檔案存在，先刪除它
			if _, err := os.Stat(newName); err == nil {
				err := os.Remove(newName)
				if err != nil {
					return fmt.Errorf("无法删除已存在的文件 %s: %w", newName, err)
				}
			}
			// 重新命名備份檔案
			err = os.Rename(path, newName)
			if err != nil {
				return fmt.Errorf("无法重命名 %s to %s: %w", path, newName, err)
			}
			totalDirs[1]++ // 已操作檔案數累加
		}
		return nil
	})

	if err != nil {
		log.Println("扫描文件夹失败:", err)
	}
	log.Printf("回滚结束。已操作文件: %d / %d\n", totalDirs[1], totalDirs[0])
}

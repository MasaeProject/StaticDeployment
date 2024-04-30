package main

import (
	"crypto/rand"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func copyFile(src string, dst string) error {
	log.Printf("复制文件: %s -> %s", src, dst) // CP
	totalIO++
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return os.ErrInvalid
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	if err := MakeDirectory(filepath.Dir(dst)); err != nil {
		return err
	}

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
}

func copyDirectory(src string, dst string) error {
	log.Printf("复制文件夹: %s -> %s", src, dst) // CPDIR
	totalIO++
	sourceInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dst, sourceInfo.Mode()); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := copyDirectory(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func CleanPath(path string) string {
	path = filepath.Clean(path)
	path = strings.ReplaceAll(path, "/", string(filepath.Separator))
	path = strings.ReplaceAll(path, "\\", string(filepath.Separator))
	return path
}

func Copy(src string, dst string) error {
	src = CleanPath(src)
	dst = CleanPath(dst)

	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if srcInfo.IsDir() {
		if dst[len(dst)-1] == filepath.Separator {
			dst = filepath.Join(dst, filepath.Base(src))
		}
		return copyDirectory(src, dst)
	} else {
		if dst[len(dst)-1] == filepath.Separator {
			dst = filepath.Join(dst, filepath.Base(src))
		}
		return copyFile(src, dst)
	}
}

func moveFile(src string, dst string) error {
	log.Printf("移动文件: %s -> %s", src, dst) // MV
	totalIO++
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dst), sourceFileStat.Mode()); err != nil {
		return err
	}

	return os.Rename(src, dst)
}

func moveDirectory(src string, dst string) error {
	log.Printf("移动文件夹: %s -> %s", src, dst) // MVDIR
	totalIO++
	sourceInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dst, sourceInfo.Mode()); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if err := Move(srcPath, dstPath); err != nil {
			return err
		}
	}

	return os.Remove(src)
}

func Move(src string, dst string) error {
	src = CleanPath(src)
	dst = CleanPath(dst)

	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if srcInfo.IsDir() {
		if dst[len(dst)-1] == filepath.Separator {
			dst = filepath.Join(dst, filepath.Base(src))
		}
		return moveDirectory(src, dst)
	} else {
		if dst[len(dst)-1] == filepath.Separator {
			dst = filepath.Join(dst, filepath.Base(src))
		}
		return moveFile(src, dst)
	}
}

func moveFileSecure(src string, dst string) error {
	log.Printf("移动文件(复制): %s -> %s", src, dst) // SMV
	totalIO++
	if err := copyFile(src, dst); err != nil {
		return err
	}
	return removeFile(src)
}

func moveDirectorySecure(src string, dst string) error {
	log.Printf("移动文件夹(复制): %s -> %s", src, dst) // SMVDIR
	totalIO++
	if err := copyDirectory(src, dst); err != nil {
		return err
	}
	return removeDirectory(src)
}

func MoveSecure(src string, dst string) error {
	src = CleanPath(src)
	dst = CleanPath(dst)

	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if srcInfo.IsDir() {
		if dst[len(dst)-1] == filepath.Separator {
			dst = filepath.Join(dst, filepath.Base(src))
		}
		return moveDirectorySecure(src, dst)
	} else {
		if dst[len(dst)-1] == filepath.Separator {
			dst = filepath.Join(dst, filepath.Base(src))
		}
		return moveFileSecure(src, dst)
	}
}

func removeFile(filePath string) error {
	log.Printf("删除文件: %s", filePath) // RM
	totalIO++
	err := os.Remove(filePath)
	if err != nil {
		return err
	}
	return nil
}

func removeDirectory(dirPath string) error {
	log.Printf("删除文件夹: %s", dirPath) // RMDIR
	totalIO++
	err := os.RemoveAll(dirPath)
	if err != nil {
		return err
	}
	return nil
}

func Remove(path string) error {
	path = CleanPath(path)

	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return removeDirectory(path)
	} else {
		return removeFile(path)
	}
}

func RemoveFileSecure(filename string) error {
	log.Printf("安全删除文件: %s", filename) // SRM
	totalIO++
	file, err := os.OpenFile(filename, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	_, err = io.CopyN(file, rand.Reader, stat.Size())
	if err != nil {
		return err
	}

	err = file.Truncate(0)
	if err != nil {
		return err
	}

	if err := file.Close(); err != nil {
		return err
	}

	return os.Remove(filename)
}

func RemoveDirectorySecure(dirPath string) error {
	log.Printf("安全删除文件夹: %s", dirPath) // SRMDIR
	totalIO++
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			if err := RemoveFileSecure(path); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return os.Remove(dirPath)
}

func RemoveSecure(path string) error {
	path = CleanPath(path)

	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return RemoveDirectorySecure(path)
	} else {
		return RemoveFileSecure(path)
	}
}

func RenamePath(oldPath, newPath string) error {
	log.Printf("重命名: %s -> %s", oldPath, newPath) // REN
	totalIO++
	err := os.Rename(oldPath, newPath)
	if err != nil {
		return err
	}
	return nil
}

func IsDirectory(path string) int8 {
	info, err := os.Stat(path)
	if err != nil {
		return -1
	}
	if info.IsDir() {
		return 1
	}
	return 0
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

func MakeDirectory(path string) error {
	var pathPart string = path
	for {
		if Exists(pathPart) {
			break
		}
		var parentDir string = filepath.Dir(pathPart)
		if parentDir == pathPart {
			break
		}
		pathPart = parentDir
	}
	sourceFileStat, err := os.Stat(pathPart)
	if err != nil {
		return err
	}
	return os.MkdirAll(path, sourceFileStat.Mode())
}

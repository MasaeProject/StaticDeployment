package main

import (
	"crypto/rand"
	"io"
	"log"
	"os"
	"path/filepath"
)

func copyFile(src, dst string) error {
	log.Printf("CP %s -> %s", src, dst)
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

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
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
	log.Printf("CPDIR %s -> %s", src, dst)
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

func Copy(src, dst string) error {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if srcInfo.IsDir() {
		if dst[len(dst)-1] == '/' {
			dst = filepath.Join(dst, filepath.Base(src))
		}
		return copyDirectory(src, dst)
	} else {
		if dst[len(dst)-1] == '/' {
			dst = filepath.Join(dst, filepath.Base(src))
		}
		return copyFile(src, dst)
	}
}

func moveFile(src, dst string) error {
	log.Printf("MV %s -> %s", src, dst)
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	return os.Rename(src, dst)
}

func moveDirectory(src, dst string) error {
	log.Printf("MVDIR %s -> %s", src, dst)
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

func Move(src, dst string) error {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if srcInfo.IsDir() {
		if dst[len(dst)-1] == '/' {
			dst = filepath.Join(dst, filepath.Base(src))
		}
		return moveDirectory(src, dst)
	} else {
		if dst[len(dst)-1] == '/' {
			dst = filepath.Join(dst, filepath.Base(src))
		}
		return moveFile(src, dst)
	}
}

func moveFileSecure(src, dst string) error {
	log.Printf("SMV %s -> %s", src, dst)
	if err := copyFile(src, dst); err != nil {
		return err
	}
	return removeFile(src)
}

func moveDirectorySecure(src, dst string) error {
	log.Printf("SMVDIR %s -> %s", src, dst)
	if err := copyDirectory(src, dst); err != nil {
		return err
	}
	return removeDirectory(src)
}

func MoveSecure(src, dst string) error {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if srcInfo.IsDir() {
		if dst[len(dst)-1] == '/' {
			dst = filepath.Join(dst, filepath.Base(src))
		}
		return moveDirectorySecure(src, dst)
	} else {
		if dst[len(dst)-1] == '/' {
			dst = filepath.Join(dst, filepath.Base(src))
		}
		return moveFileSecure(src, dst)
	}
}

func removeFile(filePath string) error {
	log.Printf("RM %s", filePath)
	err := os.Remove(filePath)
	if err != nil {
		return err
	}
	return nil
}

func removeDirectory(dirPath string) error {
	log.Printf("RMDIR %s", dirPath)
	err := os.RemoveAll(dirPath)
	if err != nil {
		return err
	}
	return nil
}

func Remove(path string) error {
	path = filepath.Clean(path)

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
	log.Printf("SRM %s", filename)
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
	log.Printf("SRMDIR %s", dirPath)
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
	path = filepath.Clean(path)

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

func Rename(oldPath, newPath string) error {
	log.Printf("REN %s -> %s", oldPath, newPath)
	err := os.Rename(oldPath, newPath)
	if err != nil {
		return err
	}
	return nil
}

//go:generate goversioninfo -icon=ico/icon.ico -manifest=main.exe.manifest -arm=true
package main

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

func StaticDeployment_Join(cmd []string) ([2]int64, error) {
	var cmdLen int = len(cmd)
	var dataLen [2]int64 = [2]int64{-1, -1}
	if cmdLen <= 2 {
		// path = srcPath
		return dataLen, fmt.Errorf("NO PATH")
	}
	var destPath string = os.Args[1]

	sourceFileStat, err := os.Stat(os.Args[2])
	if err != nil {
		log.Printf("[%s] %s : %v", os.Args[2], filepath.Dir(destPath), err)
		return dataLen, err
	}
	var fileMode fs.FileMode = sourceFileStat.Mode()

	if err := os.MkdirAll(filepath.Dir(destPath), fileMode); err != nil {
		log.Printf("[%s] %s : %v", os.Args[0], filepath.Dir(destPath), err)
		return dataLen, err
	}

	destinationFile, err := os.OpenFile(destPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, fileMode)
	if err != nil {
		log.Printf("[%s] %s : %v", os.Args[0], destPath, err)
		return dataLen, err
	}
	defer destinationFile.Close()

	if destinationFileStat, err := destinationFile.Stat(); err == nil {
		dataLen[0] = destinationFileStat.Size()
	}

	var sourcePath []string = os.Args[2:]

	for _, sourcePath := range sourcePath {
		sourceFile, err := os.Open(sourcePath)
		if err != nil {
			log.Printf("[%s] %s : %v", os.Args[0], sourcePath, err)
			return dataLen, err
		}
		defer sourceFile.Close()
		_, err = io.Copy(destinationFile, sourceFile)
		if err != nil {
			log.Printf("[%s] %s -> %s : %v", os.Args[0], destPath, sourcePath, err)
			return dataLen, err
		}
	}

	if destinationFileStat, err := destinationFile.Stat(); err == nil {
		dataLen[1] = destinationFileStat.Size()
	}

	return dataLen, nil
}

func main() {
	dataLen, err := StaticDeployment_Join(os.Args)
	log.Printf("[%s] %d -> %d  E: %v", os.Args[0], dataLen[0], dataLen[1], err)
	if err != nil {
		os.Exit(1)
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Project struct {
	Name       string        `json:"name"`
	Source     string        `json:"source"`
	Additional []string      `json:"additional"`
	Exec       [][]string    `json:"exec"`
	Replace    []ReplaceItem `json:"replace"`
}

type ReplaceItem struct {
	To         string          `json:"to"`
	Replace    []ReplaceDetail `json:"replace"`
	Additional []string        `json:"additional"`
	Exec       [][]string      `json:"exec"`
}

type ReplaceDetail struct {
	Search  string `json:"search"`
	Replace string `json:"replace"`
}

func main() {
	if len(os.Args) <= 1 {
		log.Println("必须指定一个配置文件路径。")
		return
	}
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Println("打开配置文件", os.Args[1], "失败:", err)
		return
	}
	decoder := json.NewDecoder(file)
	var projects []Project
	err = decoder.Decode(&projects)
	if err != nil {
		log.Println("配置文件解析失败:", err)
		return
	}

	fmt.Println(projects)
}

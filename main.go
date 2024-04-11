package main

import (
	"encoding/json"
	"log"
	"os"
	"runtime"
)

var (
	osName string
)

func main() {
	osName = runtime.GOOS
	log.Println("StaticDeployment v0.0.1 for", osName)

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

	for i, project := range projects {
		if runProject(project) {
			log.Printf("项目 %d : %s 处理成功。\n", i, project.Name)
		} else {
			log.Printf("项目 %d : %s 处理失败。\n", i, project.Name)
		}
	}
}

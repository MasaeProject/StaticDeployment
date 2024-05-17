package main

import (
	"fmt"
	"os"
	"testing"
)

func TestStaticDeployment_Join(t *testing.T) {
	var cmd []string = make([]string, 11)
	cmd[0] = os.Args[0]
	// 创建和写入文件
	for i := 1; i <= 10; i++ {
		filename := fmt.Sprintf("%02d.test.txt", i)
		cmd[i] = filename
		file, err := os.Create(filename)
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}
		// 将数字写入文件
		_, err = fmt.Fprintf(file, "%02d", i)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			file.Close()
			return
		}
		file.Close()
	}

	t.Log(cmd)
	// dataLen, err := StaticDeployment_Join(cmd)
	// if err == nil {
	// 	var dataLenLen int = len(dataLen)
	// 	var dataLenStrArr []string = make([]string, dataLenLen)
	// 	for i, num := range dataLen {
	// 		dataLenStrArr[i] = strconv.FormatInt(num, 10)
	// 	}
	// 	var total string = dataLenStrArr[dataLenLen-1]
	// 	dataLenStrArr[0] = strings.Join(dataLenStrArr[:dataLenLen-1], " + ")
	// 	t.Log(fmt.Printf("[%s] %s = %s (E:%v)", os.Args[0], dataLenStrArr[0], total, err))
	// }

	// for i := 1; i <= 99; i++ {
	// 	filename := fmt.Sprintf("%02d.test.txt", i)
	// 	err := os.Remove(filename)
	// 	if err != nil {
	// 		fmt.Println("Error deleting file:", err)
	// 	}
	// }

}

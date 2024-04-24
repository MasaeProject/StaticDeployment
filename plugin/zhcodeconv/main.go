package zhcodeconv

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"unicode"
)

func InitWithCmd(cmd []string, srcPath string) error {
	var cmdLen int = len(cmd)
	var path string = ""
	var outPath string = ""
	var mode string = "hex"
	if cmdLen <= 1 {
		path = srcPath
	} else if cmdLen >= 2 {
		path = cmd[1]
	}
	if cmdLen >= 3 {
		mode = cmd[2]
	}
	if cmdLen >= 4 {
		outPath = cmd[3]
	}
	if len(outPath) == 0 {
		outPath = path
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var newCode []byte = replaceNonAscii(data, mode)
	err = os.WriteFile(outPath, newCode, 0644)
	if err != nil {
		return err
	}
	return nil
}

func replaceNonAscii(code []byte, mode string) []byte {
	var quotation byte = 0
	var newCode []byte = []byte{}
	var cache []byte = []byte{}
	for _, c := range code {
		if c == '"' || c == '\'' {
			if quotation == c {
				quotation = 0
			} else if quotation == 0 {
				quotation = c
			}
		}
		if quotation == 0 && (c <= 0 || c > 127) {
			cache = append(cache, c)
		} else {
			if len(cache) > 0 {
				if mode == "hex" {
					newCode = append(newCode, []byte(firstNoNumber(fmt.Sprintf("%04x", cache)))...)
				} else if mode == "md5" {
					hasher := md5.New()
					hasher.Write([]byte(cache))
					hash := hasher.Sum(nil)
					nStr := hex.EncodeToString(hash)
					newCode = append(newCode, []byte(firstNoNumber(nStr))...)
				} else {
					newCode = append(newCode, cache...)
				}
				cache = []byte{}
			}
			newCode = append(newCode, c)
		}
	}
	return newCode
}

func firstNoNumber(str string) string {
	if unicode.IsDigit(rune(str[0])) {
		str = "g" + str
	}
	return str
}

// func main() {
// 	if len(os.Args) <= 1 {
// 		log.Println("错误: 必须指定要转换的文件路径。")
// 		os.Exit(1)
// 		return
// 	}
// 	var outFile string = ""
// 	if len(os.Args) >= 3 {
// 		outFile = os.Args[2]
// 	}
// 	var mode string = "hex"
// 	if len(os.Args) >= 4 {
// 		mode = os.Args[3]
// 	}
// 	code, err := os.ReadFile(os.Args[1])
// 	if err != nil {
// 		log.Println("错误: 打开文件失败：", err)
// 		os.Exit(1)
// 		return
// 	}
// 	var newCode string = replaceNonAscii(string(code), mode)
// 	if len(outFile) == 0 {
// 		fmt.Println(newCode)
// 	} else {
// 		err = os.WriteFile(outFile, []byte(newCode), 0644)
// 		if err != nil {
// 			log.Println("错误: 写入文件失败：", err)
// 			os.Exit(1)
// 			return
// 		}
// 	}

// 	fmt.Println()
// }

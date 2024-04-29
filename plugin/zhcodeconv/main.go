package minify

import (
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"hash"
	"os"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/crypto/sha3"
)

func InitWithCmd(cmd []string, srcPath string) ([2]int, error) {
	var cmdLen int = len(cmd)
	var path string = ""
	var outPath string = ""
	var mode string = "hex"
	var reEnc int = 1
	var DataLen [2]int = [2]int{-1, -1}
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
		return DataLen, err
	}
	var modeArr []string = strings.Split(mode, "*")
	if len(modeArr) > 1 {
		var err error
		mode = modeArr[0]
		reEnc, err = strconv.Atoi(modeArr[1])
		if err != nil {
			reEnc = 1
		}
	}
	DataLen[0] = len(data)
	var newCode []byte = replaceNonAscii(data, mode, reEnc)
	DataLen[1] = len(newCode)
	err = os.WriteFile(outPath, newCode, 0644)
	if err != nil {
		return DataLen, err
	}
	return DataLen, nil
}

func replaceNonAscii(code []byte, mode string, reEnc int) []byte {
	var quotation string = ""
	var newCode []byte = []byte{}
	var cache []byte = []byte{}
	stringToTraverse := string(code)
	for i, w := 0, 0; i < len(stringToTraverse); i += w {
		runeValue, width := utf8.DecodeRuneInString(stringToTraverse[i:])
		var c string = string(runeValue)
		// fmt.Printf("%s", c)
		w = width
		if c == "\"" || c == "'" {
			if quotation == c {
				quotation = ""
			} else if quotation == "" {
				quotation = c
			}
		}
		if len(quotation) == 0 && (c[0] <= 0 || c[0] > 127) {
			cache = append(cache, []byte(c)...)
		} else {
			if len(cache) > 0 {
				var cryptoCache []byte = hashesCrypto(mode, cache)
				for i := 1; i < reEnc; i++ {
					cryptoCache = hashesCrypto(mode, cryptoCache)
				}
				newCode = append(newCode, cryptoCache...)
				cache = []byte{}
			}
			newCode = append(newCode, []byte(c)...)
		}
	}
	return newCode
}

func hashesCrypto(mode string, cache []byte) []byte {
	var newCode []byte = []byte{}
	var hasher hash.Hash
	switch mode {
	case "hex":
		return []byte(fmt.Sprintf("%x", cache))
	case "md5":
		hasher = md5.New()
	case "sha1":
		hasher = sha1.New()
	case "sha256":
		hasher = sha3.New256()
	case "sha512":
		hasher = sha3.New512()
	default:
		return cache
	}
	_, err := hasher.Write(cache)
	if err != nil {
		fmt.Println("ERROR: hasher.Write error:", err)
		return cache
	}
	var hash []byte = hasher.Sum(nil)
	var hashStr string = fmt.Sprintf("%x", hash)
	hashStr = firstNoNumberS(hashStr)
	newCode = append(newCode, []byte(hashStr)...)
	return newCode
}

func firstNoNumberS(str string) string {
	if unicode.IsDigit(rune(str[0])) {
		str = "g" + str
	}
	return str
}
func firstNoNumberC(str []byte) []byte {
	if unicode.IsDigit(rune(str[0])) {
		str = append([]byte{'g'}, str...)
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

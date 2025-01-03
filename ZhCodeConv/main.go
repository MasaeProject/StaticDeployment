//go:generate goversioninfo -icon=ico/icon.ico -manifest=main.exe.manifest -arm=true
package main

import (
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"hash"
	"hash/adler32"
	"hash/crc32"
	"hash/crc64"
	"hash/fnv"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/crypto/sha3"
)

var title string = "[代码名称转换器] "

func StaticDeployment_ZhCodeConv(cmd []string) ([2]int, error) {
	var cmdLen int = len(cmd)
	var path string = ""
	var outPath string = ""
	var mode string = "hex"
	var reEnc int = 1
	var dataLen [2]int = [2]int{-1, -1}
	var syms []string = []string{"\"", "'"}
	if cmdLen <= 1 {
		// path = srcPath
		return dataLen, fmt.Errorf(title + "请输入文件路径")
	} else if cmdLen >= 2 {
		path = cmd[1]
	}
	if cmdLen >= 3 {
		mode = cmd[2]
	}
	if cmdLen >= 4 {
		outPath = cmd[3]
	}
	if cmdLen >= 5 {
		if len(cmd[4]) == 0 || cmd[4] == "." {
			syms = []string{}
		} else {
			syms = strings.Split(cmd[4], "")
		}
	}
	// if len(outPath) == 0 {
	// 	outPath = path
	// }
	sourceFileStat, err := os.Stat(path)
	if err != nil {
		return dataLen, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return dataLen, err
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
	dataLen[0] = len(data)
	var newCode []byte = replaceNonAscii(data, mode, reEnc, syms)
	dataLen[1] = len(newCode)
	if len(outPath) == 0 {
		fmt.Println(string(newCode))
	} else {
		err = os.WriteFile(outPath, newCode, sourceFileStat.Mode())
		if err != nil {
			return dataLen, err
		}
	}
	return dataLen, nil
}

func isEnglishLetter(r rune) bool {
	return r <= 127 && unicode.IsLetter(r)
}

func replaceNonAscii(code []byte, mode string, reEnc int, syms []string) []byte {
	var quotation string = ""
	var newCode []byte = []byte{}
	var cache []byte = []byte{}
	for i := 0; i < len(code); {
		r, size := utf8.DecodeRune(code[i:])
		var c string = string(r)
		if len(syms) > 0 {
			for _, sym := range syms {
				if c == sym {
					if quotation == c {
						quotation = ""
					} else if len(quotation) == 0 {
						quotation = c
					}
					break
				}
			}
		}
		if unicode.IsLetter(r) && len(quotation) == 0 && !isEnglishLetter(r) {
			cache = append(cache, []byte(c)...) // 中文
			i += size
			continue
		}
		if len(cache) > 0 {
			newC := EncodeCrypto(cache, mode, reEnc)
			cache = []byte{}
			newCode = append(newCode, newC...)
		}
		newCode = append(newCode, []byte(c)...)

		i += size
	}
	if len(cache) > 0 {
		newC := EncodeCrypto(cache, mode, reEnc)
		newCode = append(newCode, newC...)
	}
	return newCode
}

func EncodeCrypto(cache []byte, mode string, reEnc int) []byte {
	var cryptoCache []byte = hashesCrypto(mode, cache)
	for i := 1; i < reEnc; i++ {
		cryptoCache = hashesCrypto(mode, cryptoCache)
	}
	return cryptoCache
}

func hashesCrypto(mode string, cache []byte) []byte {
	var newCode []byte = []byte{}
	var hasher hash.Hash
	switch mode {
	// 以下為非加密雜湊函式
	case "hex":
		return []byte(fmt.Sprintf("%x", cache))
	case "fnv1a":
		hasher = fnv.New32a()
	case "adler32":
		hasher = adler32.New()
	case "crc32":
		hasher = crc32.NewIEEE()
	case "crc64":
		hasher = crc64.New(crc64.MakeTable(crc64.ISO))
	// 以下為加密雜湊函式
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
		log.Println(title+"[错误] hasher.Write:", err)
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

func main() {
	dataLen, err := StaticDeployment_ZhCodeConv(os.Args)
	log.Printf(title+"%d -> %d (E:%v)", dataLen[0], dataLen[1], err)
	if err != nil {
		os.Exit(1)
	}
}

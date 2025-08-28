package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// SizeFormat 字节的单位转换 保留两位小数
func SizeFormat(size int64) string {
	if size < 1024 {
		//return strconv.FormatInt(fileSize, 10) + "B"
		return fmt.Sprintf("%.2fB", float64(size)/float64(1))
	} else if size < (1024 * 1024) {
		return fmt.Sprintf("%.2fKB", float64(size)/float64(1024))
	} else if size < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fMB", float64(size)/float64(1024*1024))
	} else if size < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fGB", float64(size)/float64(1024*1024*1024))
	} else if size < (1024 * 1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fTB", float64(size)/float64(1024*1024*1024*1024))
	} else { //if fileSize < (1024 * 1024 * 1024 * 1024 * 1024 * 1024)
		return fmt.Sprintf("%.2fEB", float64(size)/float64(1024*1024*1024*1024*1024))
	}
}

// AnyToJsonB interface{} -> json string
func AnyToJsonB(data interface{}) ([]byte, error) {
	jsonStr, err := json.Marshal(data)
	return jsonStr, err
}

// GetMD5Encode 获取Md5编码
func GetMD5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// ParsePath 解析路径，返回目录、文件名（不含后缀）、完整文件名、后缀
func ParsePath(fullPath string) (dir, nameWithoutExt, fileName, ext string) {
	// 1. 提取完整文件名（含后缀）
	fileName = filepath.Base(fullPath)

	// 2. 提取目录路径
	dir = filepath.Dir(fullPath)

	// 3. 提取文件后缀（含 "."，如 ".png"）
	ext = filepath.Ext(fullPath)

	// 4. 提取文件名（不含后缀）
	if ext != "" {
		// 从完整文件名中去掉后缀（注意：如果文件名有多个点，只去掉最后一个点及后面的内容）
		nameWithoutExt = strings.TrimSuffix(fileName, ext)
	} else {
		nameWithoutExt = fileName // 无后缀时，文件名即本身
	}

	return dir, nameWithoutExt, fileName, ext
}

func GetFileSize(filePath string) (string, error) {
	// 获取文件信息
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return "", err
	}

	// 判断是否为文件（不是目录）
	if !fileInfo.Mode().IsRegular() {
		return "", fmt.Errorf("%s 不是常规文件", filePath)
	}

	// 返回文件大小（字节数）
	return SizeFormat(fileInfo.Size()), nil
}

// 缓冲区对象池：缓存512字节的字节切片
var bufferPool = sync.Pool{
	// New函数：当池为空时，创建新的缓冲区
	New: func() interface{} {
		return make([]byte, 512) // 固定512字节，满足http.DetectContentType需求
	},
}

func DetectByStdLib(filePath string) (string, error) {
	buf := bufferPool.Get().([]byte)
	defer bufferPool.Put(buf)

	file, err := os.OpenFile(filePath, os.O_RDONLY, 0)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = file.Close()
	}()

	n, err := file.Read(buf)
	if err != nil {
		return "", err
	}

	return http.DetectContentType(buf[:n]), nil
}

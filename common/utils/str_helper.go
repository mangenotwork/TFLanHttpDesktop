package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
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

// AnyToString any -> string
func AnyToString(i interface{}) string {
	if i == nil {
		return ""
	}
	if reflect.ValueOf(i).Kind() == reflect.String {
		return i.(string)
	}
	var buf bytes.Buffer
	stringValue(reflect.ValueOf(i), 0, &buf)
	return buf.String()
}

func stringValue(v reflect.Value, indent int, buf *bytes.Buffer) {
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Struct:
		buf.WriteString("{\n")
		for i := 0; i < v.Type().NumField(); i++ {
			ft := v.Type().Field(i)
			fv := v.Field(i)
			if ft.Name[0:1] == strings.ToLower(ft.Name[0:1]) {
				continue
			}
			if (fv.Kind() == reflect.Ptr || fv.Kind() == reflect.Slice) && fv.IsNil() {
				continue
			}
			buf.WriteString(strings.Repeat(" ", indent+2))
			buf.WriteString(ft.Name + ": ")
			if tag := ft.Tag.Get("sensitive"); tag == "true" {
				buf.WriteString("<sensitive>")
			} else {
				stringValue(fv, indent+2, buf)
			}
			buf.WriteString(",\n")
		}
		buf.WriteString("\n" + strings.Repeat(" ", indent) + "}")

	case reflect.Slice:
		nl, id, id2 := "", "", ""
		if v.Len() > 3 {
			nl, id, id2 = "\n", strings.Repeat(" ", indent), strings.Repeat(" ", indent+2)
		}
		buf.WriteString("[" + nl)
		for i := 0; i < v.Len(); i++ {
			buf.WriteString(id2)
			stringValue(v.Index(i), indent+2, buf)

			if i < v.Len()-1 {
				buf.WriteString("," + nl)
			}
		}
		buf.WriteString(nl + id + "]")

	case reflect.Map:
		buf.WriteString("{\n")
		for i, k := range v.MapKeys() {
			buf.WriteString(strings.Repeat(" ", indent+2))
			buf.WriteString(k.String() + ": ")
			stringValue(v.MapIndex(k), indent+2, buf)

			if i < v.Len()-1 {
				buf.WriteString(",\n")
			}
		}
		buf.WriteString("\n" + strings.Repeat(" ", indent) + "}")

	default:
		format := "%v"
		switch v.Interface().(type) {
		case string:
			format = "%q"
		}
		_, _ = fmt.Fprintf(buf, format, v.Interface())
	}
}

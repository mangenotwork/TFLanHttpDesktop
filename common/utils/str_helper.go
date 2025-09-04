package utils

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
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

// Int64ToStr int64 -> string
func Int64ToStr(i int64) string {
	return strconv.FormatInt(i, 10)
}

// Get16MD5Encode 返回一个16位md5加密后的字符串
func Get16MD5Encode(data string) string {
	return GetMD5Encode(data)[8:24]
}

// CompressStringToBase64 将字符串压缩并返回 base64 编码的字符串
func CompressStringToBase64(s string) (string, error) {
	var buf bytes.Buffer

	// 创建 gzip 压缩 writer
	gz := gzip.NewWriter(&buf)
	_, err := gz.Write([]byte(s))
	if err != nil {
		return "", err
	}

	// 必须 Close() 才会刷新并写入 gzip 尾部
	if err := gz.Close(); err != nil {
		return "", err
	}

	// 将压缩后的二进制数据用 base64 编码为字符串
	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())
	return encoded, nil
}

// DecompressBase64ToString 将 base64 编码的压缩字符串解压回原始字符串
func DecompressBase64ToString(encoded string) ([]byte, error) {
	// 先 base64 解码
	compressedData, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}

	// 创建 gzip 读取器
	reader, err := gzip.NewReader(bytes.NewReader(compressedData))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	// 读取解压后的数据
	var buf bytes.Buffer
	_, err = io.Copy(&buf, reader)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// FileExists 判断路径（文件/目录）是否存在，并返回具体错误
// 参数：
//
//	path: 待判断的路径
//	followLink: 是否跟随符号链接（true=跟随目标，false=仅判断链接本身）
//
// 返回值：
//
//	exists: true=存在，false=不存在
//	err: 非nil表示判断过程出错（如权限不足）
func FileExists(path string, followLink bool) (exists bool, err error) {
	// 路径预处理：清理冗余字符并适配跨平台
	cleanPath := filepath.Clean(path)
	if cleanPath == "" {
		return false, errors.New("无效的空路径")
	}

	// 根据是否跟随符号链接选择不同的Stat函数
	var errStat error
	if followLink {
		_, errStat = os.Stat(cleanPath) // 跟随符号链接
	} else {
		_, errStat = os.Lstat(cleanPath) // 不跟随符号链接
	}

	// 错误处理逻辑
	switch {
	case errStat == nil:
		// 无错误 → 路径存在
		return true, nil
	case os.IsNotExist(errStat):
		// 明确不存在 → 返回false（无错误）
		return false, nil
	case os.IsPermission(errStat):
		// 权限不足 → 返回错误
		return false, errors.Join(
			errors.New("权限不足，无法访问路径"),
			errStat,
			errors.New("路径: "+cleanPath),
		)
	default:
		// 其他系统错误（如路径非法、IO错误等）
		return false, errors.Join(
			errors.New("判断路径存在性失败"),
			errStat,
			errors.New("路径: "+cleanPath),
		)
	}
}

// FileExistsDefault 简化版：默认跟随符号链接，仅返回存在性（兼容原用法）
func FileExistsDefault(path string) bool {
	exists, _ := FileExists(path, true)
	return exists
}

func SliceDeduplicate[V comparable](a []V) []V {
	l := len(a)
	if l < 2 {
		return a
	}
	seen := make(map[V]struct{})
	j := 0
	for i := 0; i < l; i++ {
		if _, ok := seen[a[i]]; ok {
			continue
		}
		seen[a[i]] = struct{}{}
		a[j] = a[i]
		j++
	}
	return a[:j]
}

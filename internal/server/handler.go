package server

import (
	"TFLanHttpDesktop/common/logger"
	"TFLanHttpDesktop/internal/server/assets"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func Health(ctx *gin.Context) {
	ctx.String(http.StatusOK, "ok")
	return
}

func DebugDownloadPg(ctx *gin.Context) {
	ctx.Data(http.StatusOK, "text/html; charset=utf-8", []byte(assets.DownloadPg))
}

func DebugUploadPg(ctx *gin.Context) {
	ctx.Data(http.StatusOK, "text/html; charset=utf-8", []byte(assets.UploadPg))
}

func DebugMemoPg(ctx *gin.Context) {
	ctx.Data(http.StatusOK, "text/html; charset=utf-8", []byte(assets.MemoPg))
}

func Tailwindcss(ctx *gin.Context) {
	logger.Info("Tailwindcss...")
	//b, err := os.ReadFile("./internal/server/assets/tailwindcss.js")
	//if err != nil {
	//	logger.Error("读取文件失败,", err.Error())
	//}
	//
	//b64, _ := compressStringToBase64(string(b))
	//logger.Info(b64)

	b, err := decompressBase64ToString(assets.TailwindcssData)
	if err != nil {
		logger.Error("读取文件失败,", err.Error())
	}
	ctx.Data(http.StatusOK, "text/javascript", b)
}

// compressStringToBase64 将字符串压缩并返回 base64 编码的字符串
func compressStringToBase64(s string) (string, error) {
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

// decompressBase64ToString 将 base64 编码的压缩字符串解压回原始字符串
func decompressBase64ToString(encoded string) ([]byte, error) {
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

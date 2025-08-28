package server

import (
	"TFLanHttpDesktop/common/define"
	"TFLanHttpDesktop/common/logger"
	"TFLanHttpDesktop/common/utils"
	"TFLanHttpDesktop/internal/server/assets"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
)

func Health(ctx *gin.Context) {
	ctx.String(http.StatusOK, "ok")
	return
}

func DebugDownloadPg(ctx *gin.Context) {
	ctx.Data(http.StatusOK, "text/html; charset=utf-8", []byte(assets.DownloadPg))
}

func DownloadPg(ctx *gin.Context) {
	fileKey := ctx.Param("file")
	logger.Debug("fileKey = ", fileKey)
	fileKey = strings.Replace(fileKey, "/", "", -1)
	logger.Debug("define.DownloadMem = ", define.DownloadMem)
	filePath, ok := define.DownloadMem[fileKey]
	if !ok {
		logger.Debug("file not found")
		ctx.Data(http.StatusOK, "text/html; charset=utf-8", []byte("下载链接已失效"))
		return
	}
	logger.Info(filePath)

	_, _, fileName, ext := utils.ParsePath(filePath)
	fileSize, _ := utils.GetFileSize(filePath)

	tpl, err := template.New("html").Parse(assets.DownloadPg)
	if err != nil {
		logger.Error(err)
		ctx.Data(http.StatusInternalServerError, "text/html; charset=utf-8", []byte(err.Error()))
		return
	}
	var renderedHTML strings.Builder
	data := map[string]interface{}{
		"Title":       "下载文件",
		"FileName":    fileName,
		"Ext":         ext,
		"FileSize":    fileSize,
		"DownloadUrl": fmt.Sprintf("%s/d/%s", define.DoMain, fileKey),
	}
	if err := tpl.Execute(&renderedHTML, data); err != nil {
		logger.Error(err)
		ctx.Data(http.StatusInternalServerError, "text/html; charset=utf-8", []byte(err.Error()))
		return
	}

	ctx.Data(http.StatusOK, "text/html; charset=utf-8", []byte(renderedHTML.String()))
	return
}

func DownloadExecute(ctx *gin.Context) {
	fileKey := ctx.Param("file")
	logger.Debug("fileKey = ", fileKey)
	fileKey = strings.Replace(fileKey, "/", "", -1)
	logger.Debug("define.DownloadMem = ", define.DownloadMem)
	filePath, ok := define.DownloadMem[fileKey]
	if !ok {
		logger.Debug("file not found")
		ctx.Data(http.StatusOK, "text/html; charset=utf-8", []byte("下载链接已失效"))
		return
	}
	logger.Info(filePath)

	fileName := filepath.Base(filePath)
	encodedFileName := url.QueryEscape(fileName)
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s;", encodedFileName))
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.File(filePath)
	return
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

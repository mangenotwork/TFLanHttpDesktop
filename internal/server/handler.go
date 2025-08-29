package server

import (
	"TFLanHttpDesktop/common/define"
	"TFLanHttpDesktop/common/logger"
	"TFLanHttpDesktop/common/utils"
	"TFLanHttpDesktop/internal/data"
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
	"os"
	"path/filepath"
	"strings"
	"time"
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
		ctx.Data(http.StatusForbidden, "text/html; charset=utf-8", []byte("下载链接已失效"))
		return
	}
	logger.Info(filePath)

	downloadData, _ := data.GetDownloadData()
	logger.Debug("downloadData = ", downloadData)
	if downloadData.Path != filePath {
		ctx.Data(http.StatusForbidden, "text/html; charset=utf-8", []byte("下载链接已失效"))
		return
	}

	isPassword := 0
	if downloadData.IsPassword {
		isPassword = 1
	}

	_, _, fileName, _ := utils.ParsePath(filePath)
	fileSize, _ := utils.GetFileSize(filePath)

	tpl, err := template.New("html").Parse(assets.DownloadPg)
	if err != nil {
		logger.Error(err)
		ctx.Data(http.StatusInternalServerError, "text/html; charset=utf-8", []byte(err.Error()))
		return
	}
	var renderedHTML strings.Builder
	values := map[string]interface{}{
		"Title":       "下载-" + fileName,
		"FileName":    fileName,
		"FileSize":    fileSize,
		"DownloadUrl": fmt.Sprintf("%s/d/%s", define.DoMain, fileKey),
		"IsPassword":  isPassword,
	}
	if err := tpl.Execute(&renderedHTML, values); err != nil {
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
		ctx.Data(http.StatusForbidden, "text/html; charset=utf-8", []byte("下载链接已失效"))
		return
	}
	logger.Info(filePath)

	password := ctx.Query("p")
	logger.InfoF("password = %s", password)

	downloadData, _ := data.GetDownloadData()
	if downloadData.Path != filePath {
		ctx.Data(http.StatusForbidden, "text/html; charset=utf-8", []byte("下载链接已失效"))
		return
	}
	if downloadData.IsPassword && downloadData.Password != password {
		ctx.Data(http.StatusForbidden, "text/html; charset=utf-8", []byte("密码错误"))
		return
	}

	ua := ctx.Request.UserAgent()
	ip, _ := ctx.Get(ReqIP)
	logger.Debug("ua = ", ua)
	logger.Debug("ip = ", ip)
	fileSize, _ := utils.GetFileSize(filePath)
	err := data.SetDownloadLog(&data.DownloadLog{
		Time:      time.Now().Format(utils.TimeTemplate),
		IP:        ip.(string),
		UserAgent: ua,
		Path:      filePath,
		Size:      fileSize,
	})
	if err != nil {
		logger.Error("记录下载日志出现错误 ", err)
	}

	fileName := filepath.Base(filePath)
	encodedFileName := url.QueryEscape(fileName)
	ctx.Header("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	ctx.Header("Pragma", "no-cache") // 兼容HTTP/1.0
	ctx.Header("Expires", "0")       // 告诉浏览器该资源已过期
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s;", encodedFileName))
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.File(filePath)
	return
}

func DebugUploadPg(ctx *gin.Context) {
	ctx.Data(http.StatusOK, "text/html; charset=utf-8", []byte(assets.UploadPg))
}

func UploadPg(ctx *gin.Context) {
	fileKey := ctx.Param("file")
	logger.Debug("fileKey = ", fileKey)
	fileKey = strings.Replace(fileKey, "/", "", -1)
	logger.Debug("define.DownloadMem = ", define.UploadMem)
	filePath, ok := define.UploadMem[fileKey]
	if !ok {
		logger.Debug("file not found")
		ctx.Data(http.StatusForbidden, "text/html; charset=utf-8", []byte("下载链接已失效"))
		return
	}
	logger.Info(filePath)

	uploadData, _ := data.GetUploadData()
	logger.Debug("uploadData = ", uploadData)
	if uploadData.Path != filePath {
		ctx.Data(http.StatusForbidden, "text/html; charset=utf-8", []byte("下载链接已失效"))
		return
	}

	isPassword := 0
	if uploadData.IsPassword {
		isPassword = 1
	}

	tpl, err := template.New("html").Parse(assets.UploadPg)
	if err != nil {
		logger.Error(err)
		ctx.Data(http.StatusInternalServerError, "text/html; charset=utf-8", []byte(err.Error()))
		return
	}

	token, _ := utils.GenerateSignature(filePath)
	logger.Debug(token)

	var renderedHTML strings.Builder
	values := map[string]interface{}{
		"Title":      "上传文件",
		"UploadUrl":  fmt.Sprintf("%s/u/%s", define.DoMain, fileKey),
		"IsPassword": isPassword,
		"Token":      token,
	}
	if err := tpl.Execute(&renderedHTML, values); err != nil {
		logger.Error(err)
		ctx.Data(http.StatusInternalServerError, "text/html; charset=utf-8", []byte(err.Error()))
		return
	}

	ctx.Data(http.StatusOK, "text/html; charset=utf-8", []byte(renderedHTML.String()))
	return
}

func UploadExecute(ctx *gin.Context) {
	fileKey := ctx.Param("file")
	logger.Debug("fileKey = ", fileKey)
	fileKey = strings.Replace(fileKey, "/", "", -1)
	logger.Debug("define.DownloadMem = ", define.UploadMem)
	filePath, ok := define.UploadMem[fileKey]
	if !ok {
		logger.Debug("file not found")
		ctx.Data(http.StatusForbidden, "text/html; charset=utf-8", []byte("下载链接已失效"))
		return
	}
	logger.Info(filePath)

	uploadData, _ := data.GetUploadData()
	logger.Debug("uploadData = ", uploadData)
	if uploadData.Path != filePath {
		ctx.Data(http.StatusForbidden, "text/html; charset=utf-8", []byte("下载链接已失效"))
		return
	}

	token := ctx.PostForm("token")
	logger.Debug("token = ", token)

	if verify, _ := utils.VerifySignature(filePath, token); !verify {
		ctx.Data(http.StatusForbidden, "text/html; charset=utf-8", []byte("token无效"))
		return
	}

	password := ctx.PostForm("password")
	logger.Debug("password = ", password)
	if uploadData.IsPassword && uploadData.Password != password {
		ctx.JSON(http.StatusForbidden, "密码错误")
		return
	}

	fromData, err := ctx.MultipartForm()
	if err != nil {
		logger.Error("获取参数失败 err = ", err)
	}
	logger.Info("获取参数 fromData = ", fromData)

	files := fromData.File["files"]
	logger.Debug(files)
	if len(files) == 0 {
		ctx.JSON(http.StatusForbidden, "未上传任何文件")
		return
	}

	saveErr := make([]string, 0)
	_ = os.MkdirAll(uploadData.Path, 0755) // 确保目录存在
	for i, file := range files {
		// 构建保存路径
		dst := fmt.Sprintf("%s/%s", uploadData.Path, file.Filename)
		// todo... 判断是否已经存在，存在则从命名
		// 保存文件
		if err := ctx.SaveUploadedFile(file, dst); err != nil {
			logger.Error("保存文件失败: ", err)
			saveErr = append(saveErr, fmt.Sprintf("保存失败:%s", file.Filename))
		}
		logger.DebugF("文件 %d 保存成功: %s\n", i+1, dst)
	}

	if len(saveErr) > 0 {
		ctx.JSON(http.StatusForbidden, saveErr)
		return
	}

	ctx.JSON(http.StatusOK, "保存成功")
	return

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

	//b, err := decompressBase64ToString(assets.TailwindcssData)
	//if err != nil {
	//	logger.Error("读取文件失败,", err.Error())
	//}
	ctx.Data(http.StatusOK, "text/javascript", []byte(assets.TailwindcssData))
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

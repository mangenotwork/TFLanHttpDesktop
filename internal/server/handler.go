package server

import (
	"TFLanHttpDesktop/common/define"
	"TFLanHttpDesktop/common/logger"
	"TFLanHttpDesktop/common/utils"
	"TFLanHttpDesktop/internal/data"
	"TFLanHttpDesktop/internal/mq"
	"TFLanHttpDesktop/internal/server/assets"
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
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

	// 通知ui界面更新
	mq.Producer(&mq.ChanData{
		Type: 2,
		Msg:  fmt.Sprintf("三方设备下载了文件（ip:%s）: %s", ip, filePath),
	})

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
		ctx.Data(http.StatusUnauthorized, "text/html; charset=utf-8", []byte("密码错误"))
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
	filesName := make([]string, 0)
	for i, file := range files {
		// 构建保存路径
		dst := fmt.Sprintf("%s/%s", uploadData.Path, file.Filename)
		if utils.FileExistsDefault(dst) {
			_, name, _, ext := utils.ParsePath(file.Filename)
			newName := name + "_" + time.Now().Format(utils.TimeNumberTemplate) + ext
			dst = fmt.Sprintf("%s/%s", uploadData.Path, newName)
		}
		// 保存文件
		if err := ctx.SaveUploadedFile(file, dst); err != nil {
			logger.Error("保存文件失败: ", err)
			saveErr = append(saveErr, fmt.Sprintf("保存失败:%s", file.Filename))
		}
		logger.DebugF("文件 %d 保存成功: %s\n", i+1, dst)
		filesName = append(filesName, file.Filename)
	}

	if len(saveErr) > 0 {
		ctx.JSON(http.StatusForbidden, saveErr)
		return
	}
	ua := ctx.Request.UserAgent()
	ip, _ := ctx.Get(ReqIP)
	_ = data.SetUploadLog(&data.UploadLog{
		Time:      time.Now().Format(utils.TimeTemplate),
		IP:        ip.(string),
		UserAgent: ua,
		Path:      filePath,
		Files:     strings.Join(filesName, ","),
	})

	// 通知ui界面更新
	mq.Producer(&mq.ChanData{
		Type: 2,
		Msg:  fmt.Sprintf("三方设备上传了文件（ip:%s|路径:%s）: %s", ip, filePath, strings.Join(filesName, ",")),
	})

	ctx.JSON(http.StatusOK, "保存成功")
	return

}

func DebugMemoPg(ctx *gin.Context) {
	ctx.Data(http.StatusOK, "text/html; charset=utf-8", []byte(assets.MemoPg))
}

func Tailwindcss(ctx *gin.Context) {
	ctx.Data(http.StatusOK, "text/javascript", []byte(assets.TailwindcssData))
}

func MemoPg(ctx *gin.Context) {
	id := ctx.Param("id")
	logger.Info("id = ", id)

	memoData, err := data.GetMemoInfo(id)
	if err != nil {
		ctx.Data(http.StatusForbidden, "text/html; charset=utf-8", []byte("获取备忘录信息错误"))
		return
	}

	// 没有权限
	if memoData.Authority == 1 {
		ctx.Data(http.StatusForbidden, "text/html; charset=utf-8", []byte("没有权限访问"))
		return
	}

	memoContent, err := data.GetMemoContent(id)
	if err != nil {
		ctx.Data(http.StatusForbidden, "text/html; charset=utf-8", []byte("获取备忘录内容错误"))
		return
	}

	tpl, err := template.New("html").Parse(assets.MemoPg)
	if err != nil {
		logger.Error(err)
		ctx.Data(http.StatusInternalServerError, "text/html; charset=utf-8", []byte(err.Error()))
		return
	}
	sign, _ := ctx.Cookie("sign")
	signOk, _ := utils.VerifySignature(id+memoData.Password, sign)
	isPassword := 0
	if len(memoData.Password) > 0 && !signOk {
		isPassword = 1
	}

	var renderedHTML strings.Builder
	values := map[string]interface{}{
		"Title":      memoData.Name,
		"Content":    memoContent,
		"Authority":  memoData.Authority,
		"IsPassword": isPassword,
		"Id":         memoData.Id,
		"SaveUrl":    "/memo_save/" + memoData.Id,
	}
	if err := tpl.Execute(&renderedHTML, values); err != nil {
		logger.Error(err)
		ctx.Data(http.StatusInternalServerError, "text/html; charset=utf-8", []byte(err.Error()))
		return
	}

	ctx.Data(http.StatusOK, "text/html; charset=utf-8", []byte(renderedHTML.String()))
	return
}

func MemoVerify(ctx *gin.Context) {
	id := ctx.Query("id")
	password := ctx.PostForm("password")
	logger.Debug("id = ", id)
	logger.Debug("password = ", password)

	memoData, err := data.GetMemoInfo(id)
	if err != nil {
		ctx.Data(http.StatusForbidden, "text/html; charset=utf-8", []byte("获取备忘录信息错误"))
		return
	}

	if memoData.Password != password {
		ctx.Data(http.StatusForbidden, "text/html; charset=utf-8", []byte("密码错误"))
		return
	}

	sign, _ := utils.GenerateSignature(id + password)
	ctx.SetCookie("sign", sign, 3600*24, "/", "", false, true)
	ctx.Redirect(http.StatusFound, "/memo/"+id)

}

func MemoSave(ctx *gin.Context) {
	id := ctx.Param("id")
	logger.Info("id = ", id)

	memoData, err := data.GetMemoInfo(id)
	if err != nil {
		ctx.Data(http.StatusForbidden, "text/html; charset=utf-8", []byte("获取备忘录信息错误"))
		return
	}

	sign, _ := ctx.Cookie("sign")

	logger.Debug("sign = ", sign)

	signOk, _ := utils.VerifySignature(id+memoData.Password, sign)
	if len(memoData.Password) > 0 && !signOk {
		ctx.Data(http.StatusUnauthorized, "text/html; charset=utf-8", []byte("签名错误"))
		return
	}

	content := ctx.PostForm("content")
	_, err = data.SetMemoContent(id, content)
	if err != nil {
		logger.Error(err)
		ctx.Data(http.StatusForbidden, "text/html; charset=utf-8", []byte("保存失败"))
		return
	}

	// 通知ui界面更新
	ip, _ := ctx.Get(ReqIP)
	mq.Producer(&mq.ChanData{
		Type: 3,
		Msg:  "三方设备修改了备忘录（ip: " + ip.(string) + "）: " + memoData.Name,
	})

	ctx.JSON(http.StatusOK, "上传成功")
	return
}

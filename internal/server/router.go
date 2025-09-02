package server

import (
	"TFLanHttpDesktop/common/logger"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"strings"
	"time"
)

var Router *gin.Engine

func Routers() *gin.Engine {

	Router.Use(gzip.Gzip(gzip.DefaultCompression), CorsHandler(), Base())

	Router.GET("/health", Health)

	Router.GET("/debug/download", DebugDownloadPg)
	Router.GET("/debug/upload", DebugUploadPg)
	Router.GET("/debug/memo", DebugMemoPg)

	Router.GET("/tailwindcss", Tailwindcss)

	Router.GET("/download/*file", DownloadPg)
	Router.GET("/d/*file", DownloadExecute)

	Router.GET("/upload/*file", UploadPg)
	Router.POST("/u/*file", UploadExecute)

	Router.GET("/memo/:id", MemoPg)

	return Router
}

// CorsHandler 跨域中间件
func CorsHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "*")

		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Methods", "*")
		ctx.Header("Access-Control-Allow-Headers", "*")
		ctx.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, "+
			"Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,"+
			"Pragma,FooBar") // 跨域关键设置 让浏览器可以解析
		ctx.Header("Access-Control-Allow-Credentials", "false") //  跨域请求是否需要带cookie信息 默认设置为true

		if ctx.Request.Method == http.MethodOptions {
			ctx.JSON(http.StatusOK, "Options Request!")
		}
		ctx.Next()
	}
}

func Base() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		startTime := time.Now()

		// 设置请求端ip
		setIP(ctx)

		ctx.Next()

		reqLog(ctx, startTime)

	}
}

func reqLog(ctx *gin.Context, startTime time.Time) {
	endTime := time.Now()
	latencyTime := endTime.Sub(startTime)
	reqMethod := ctx.Request.Method
	reqUri := ctx.Request.RequestURI
	statusCode := ctx.Writer.Status()
	clientIP := ctx.ClientIP()

	logger.InfoF(" %3d | %13v | %15s | %s | %s ",
		statusCode,
		latencyTime,
		clientIP,
		reqMethod,
		reqUri)
}

const ReqIP = "ReqIP"

func setIP(ctx *gin.Context) {
	ctx.Set(ReqIP, GetIP(ctx.Request))
}

func GetIP(r *http.Request) (ip string) {
	for _, ip := range strings.Split(r.Header.Get("X-Forward-For"), ",") {
		if net.ParseIP(ip) != nil {
			return ip
		}
	}
	if ip = r.Header.Get("X-Real-IP"); net.ParseIP(ip) != nil {
		return ip
	}
	if ip, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		if net.ParseIP(ip) != nil {
			return ip
		}
	}
	return "0.0.0.0"
}

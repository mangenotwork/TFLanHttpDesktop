package server

import (
	"TFLanHttpDesktop/common/define"
	"TFLanHttpDesktop/common/logger"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"
)

func InitHttpServer() {
	go func() {
		listener, err := net.Listen("tcp", "0.0.0.0:0") // 关键：绑定0.0.0.0确保外部可访问
		if err != nil {
			logger.ErrorF("创建监听器失败: %s", err.Error())
			return
		}

		addr := listener.Addr().(*net.TCPAddr)
		actualPort := addr.Port

		lanIp, _ := GetLocalIP()
		logger.Info("局域网ip ", lanIp)

		define.HttpPort = actualPort
		define.LanIP = lanIp
		define.DoMain = fmt.Sprintf("http://%s:%d", define.LanIP, define.HttpPort)
		logger.InfoF("http服务启动 %s/health", define.DoMain)

		srv := &http.Server{
			Handler:        Routers(),
			ReadTimeout:    90 * time.Second,
			WriteTimeout:   90 * time.Second,
			MaxHeaderBytes: 1 << 21,
		}

		if err := srv.Serve(listener); err != nil && errors.Is(err, http.ErrServerClosed) {
			logger.ErrorF("http服务出现异常:%s\n", err.Error())
		}

	}()
}

func GetLocalIP() (string, error) {
	// 获取所有网络接口
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", fmt.Errorf("获取网络接口失败: %v", err)
	}

	// 遍历所有接口
	for _, iface := range interfaces {
		// 跳过禁用的接口
		if iface.Flags&net.FlagUp == 0 {
			continue
		}
		// 跳过回环接口
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		// 获取接口的所有地址
		addrs, err := iface.Addrs()
		if err != nil {
			return "", fmt.Errorf("获取接口地址失败: %v", err)
		}

		// 遍历地址，筛选 IPv4 地址
		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok {
				continue // 不是 IP 网络地址
			}
			// 筛选 IPv4 且非回环地址
			if ipNet.IP.To4() != nil && !ipNet.IP.IsLoopback() {
				return ipNet.IP.String(), nil
			}
		}
	}

	return "", fmt.Errorf("未找到有效的局域网 IPv4 地址")
}

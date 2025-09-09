package server

import (
	"TFLanHttpDesktop/common/logger"
	"errors"
	"net"
	"net/http"
	"time"
)

func InitHttpServer(listener net.Listener) {
	go func() {

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

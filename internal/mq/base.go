package mq

import (
	"TFLanHttpDesktop/common/logger"
	"fyne.io/fyne/v2"
)

type ChanData struct {
	Type int // 1 上传 2 下载 3 更新备忘录
	Msg  string
}

var Chan = make(chan *ChanData)

func Producer(data *ChanData) {
	Chan <- data
}

func RunMq() {
	go func() {
		for {
			select {
			case c := <-Chan:
				logger.Debug("消费 = ", c)

				fyne.CurrentApp().SendNotification(&fyne.Notification{
					Title:   "TFLanHttpDesktop",
					Content: c.Msg,
				})
			}
		}

	}()
}

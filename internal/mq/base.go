package mq

import (
	"TFLanHttpDesktop/common/logger"
	"TFLanHttpDesktop/common/utils"
	"TFLanHttpDesktop/internal/data"
	"TFLanHttpDesktop/internal/ui"
	"fyne.io/fyne/v2"
	"time"
)

type ChanData struct {
	Type   int // 1 上传 2 下载 3 更新备忘录
	Msg    string
	MemoId string
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

				_ = data.SetOperationLog(&data.OperationLog{
					Time:  time.Now().Format(utils.TimeTemplate),
					Event: c.Msg,
				})

				if c.Type == 3 && ui.NowMemoId == c.MemoId {
					newContent, err := data.GetMemoContent(c.MemoId)
					if err != nil {
						logger.Error(err)
					} else {
						fyne.Do(func() {
							ui.MemoEntry.SetText(newContent.String())
						})
					}

				}

			}
		}

	}()
}

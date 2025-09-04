package ui

import (
	"TFLanHttpDesktop/common/logger"
	"TFLanHttpDesktop/common/utils"
	"TFLanHttpDesktop/internal/data"
	"bytes"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"log"
	"os"
	"time"
)

var authorityMap = map[int]string{
	1: "无权限",
	2: "只读",
	3: "可读写",
}

func NewMemoEvent(isEdit bool, memoId string) {
	var err error
	oldMemoData := &data.Memo{}
	if isEdit {
		oldMemoData, err = data.GetMemoInfo(memoId)
		if err != nil {
			dialog.ShowError(err, MainWindow)
			return
		}
	}

	name := widget.NewEntry()
	password := widget.NewPasswordEntry()
	authorityValue := 3
	authority := widget.NewRadioGroup([]string{"无权限", "只读", "可读写"}, func(value string) {
		logger.Debug(value)
		switch value {
		case "无权限":
			authorityValue = 1
		case "只读":
			authorityValue = 2
		case "可读写":
			authorityValue = 3
		}
	})
	authority.Horizontal = true
	authority.SetSelected("可读写")
	authority.Required = true
	items := []*widget.FormItem{
		{Text: "标题", Widget: name, HintText: "标题，非必填"},
		{Text: "权限", Widget: authority, HintText: "该权限只针对三方设备"},
		{Text: "密码", Widget: password, HintText: "密码，非必填"},
	}

	dialogTitle := "新建备忘录"
	dialogConfirm := "创建"
	if isEdit {
		dialogTitle = fmt.Sprintf("编辑 - %s", oldMemoData.Name)
		dialogConfirm = "保存编辑"
		name.SetText(oldMemoData.Name)
		password.SetText(oldMemoData.Password)
		authority.SetSelected(authorityMap[oldMemoData.Authority])
	}

	passwordDialog := dialog.NewForm(dialogTitle, dialogConfirm, "取消", items, func(b bool) {

		logger.Debug("name = ", name.Text)
		logger.Debug("authorityValue = ", authorityValue)
		logger.Debug("password = ", password.Text)

		if len(name.Text) == 0 {
			name.Text = time.Now().Format(utils.TimeTemplate)
		}

		if isEdit {
			// 编辑
			_, err := data.SetMemoInfo(memoId, name.Text, authorityValue, password.Text)
			if err != nil {
				logger.Error(err)
				dialog.ShowError(err, MainWindow)
				return
			}
			_ = data.SetOperationLog(&data.OperationLog{
				Time:  time.Now().Format(utils.TimeTemplate),
				Event: fmt.Sprintf("编辑了备忘录:%s -> %s ", oldMemoData.Name, name.Text),
			})
		} else {
			// 新建
			_, err := data.NewMemo(name.Text, authorityValue, password.Text)
			if err != nil {
				logger.Error(err)
				dialog.ShowError(err, MainWindow)
				return
			}
			_ = data.SetOperationLog(&data.OperationLog{
				Time:  time.Now().Format(utils.TimeTemplate),
				Event: fmt.Sprintf("新建了备忘录: %s", name.Text),
			})
		}

		MemoListShow()

	}, MainWindow)

	passwordDialog.Resize(fyne.NewSize(500, 300))
	passwordDialog.Show()
}

func ImportTxtEvent() {
	fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, MainWindow)
			return
		}
		if reader == nil {
			logger.Debug("Cancelled")
			return
		}

		if reader == nil {
			logger.Debug("Cancelled")
			return
		}
		defer reader.Close()

		logger.Debug(reader.URI().Path())

		content, err := os.ReadFile(reader.URI().Path())
		if err != nil {
			logger.Error(err)
			dialog.ShowError(err, MainWindow)
			return
		}

		name := fmt.Sprintf("%s - %s - 导入", reader.URI().Name(), time.Now().Format(utils.TimeTemplate))
		memoData, err := data.NewMemo(name, 3, "")
		if err != nil {
			logger.Error(err)
			dialog.ShowError(err, MainWindow)
			return
		}
		_, err = data.SetMemoContent(memoData.Id, string(content))
		if err != nil {
			logger.Error(err)
			dialog.ShowError(err, MainWindow)
			return
		}
		_ = data.SetOperationLog(&data.OperationLog{
			Time:  time.Now().Format(utils.TimeTemplate),
			Event: fmt.Sprintf("导入了备忘录:%s ", name),
		})
		MemoListShow()

		return

	}, MainWindow)
	fd.SetFilter(storage.NewExtensionFileFilter([]string{".txt"}))
	fd.Show()
}

func CopyMemoEvent(memoUrl string) {
	clipboard := MainApp.Clipboard()
	clipboard.SetContent(memoUrl)
	dialog.ShowInformation("复制成功", "链接已复制到剪贴板!", MainWindow)
}

func OpenMemoEvent(memoUrl string) {
	qrImg, _ := utils.GetQRCodeIO(memoUrl)
	reader := bytes.NewReader(qrImg)
	DownloadQr := canvas.NewImageFromReader(reader, "移动设备在同一WiFi内扫码下载")
	DownloadQr.FillMode = canvas.ImageFillOriginal
	qrDialog := dialog.NewCustom("扫码访问", "关闭", container.NewCenter(DownloadQr), MainWindow)
	qrDialog.Resize(fyne.NewSize(500, 600))
	qrDialog.Show()
}

func DelMemoEvent() {
	dialog.ShowConfirm("确认删除", "确认删除吗?", func(b bool) {
		logger.Debug(b)
		if b {
			logger.Debug("删除 ", NowMemoId)
			err := data.DeleteMemo(NowMemoId)
			if err != nil {
				dialog.ShowError(err, MainWindow)
				return
			}
			oldMemoData := &data.Memo{}
			oldMemoData, err = data.GetMemoInfo(NowMemoId)
			if err != nil {
				dialog.ShowError(err, MainWindow)
				return
			}
			_ = data.SetOperationLog(&data.OperationLog{
				Time:  time.Now().Format(utils.TimeTemplate),
				Event: fmt.Sprintf("删除了备忘录:%s ", oldMemoData.Name),
			})
			MemoListShow()
			MemoEntryContainer.RemoveAll()
			MemoEntryContainer.Refresh()
		}
	}, MainWindow)
}

func MemoSaveToTxt() {
	memoData, err := data.GetMemoInfo(NowMemoId)
	if err != nil {
		dialog.ShowError(err, MainWindow)
		return
	}

	fd := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			dialog.ShowError(err, MainWindow)
			return
		}
		if writer == nil {
			log.Println("Cancelled")
			return
		}

		file, err := os.Create(writer.URI().Path())
		if err != nil {
			fmt.Println("创建文件失败：", err)
			return
		}
		defer file.Close()

		_, err = file.WriteString(MemoEntry.Text)
		if err != nil {
			logger.Error("写入文件失败：", err)
			dialog.ShowError(err, MainWindow)
			return
		}
		logger.Info("另存为成功")
		oldMemoData := &data.Memo{}
		oldMemoData, err = data.GetMemoInfo(NowMemoId)
		if err != nil {
			dialog.ShowError(err, MainWindow)
			return
		}
		_ = data.SetOperationLog(&data.OperationLog{
			Time:  time.Now().Format(utils.TimeTemplate),
			Event: fmt.Sprintf("另存了备忘录:%s -> %s ", oldMemoData.Name, writer.URI().Path()),
		})
		dialog.ShowInformation("另存成功", fmt.Sprintf("另存至:\n%s", writer.URI().Path()), MainWindow)
	}, MainWindow)
	fd.SetFilter(storage.NewExtensionFileFilter([]string{".txt"}))
	fd.SetFileName(memoData.Name + ".txt")
	fd.SetTitleText("另存为txt")
	fd.Show()
}

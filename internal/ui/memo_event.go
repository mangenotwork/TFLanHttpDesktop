package ui

import (
	"TFLanHttpDesktop/common/define"
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
	1: MLGet(MLTNoPermission),
	2: MLGet(MLTReadOnly),
	3: MLGet(MLTReadWrite),
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
	authority := widget.NewRadioGroup([]string{MLGet(MLTNoPermission), MLGet(MLTReadOnly), MLGet(MLTReadWrite)}, func(value string) {
		logger.Debug(value)
		switch value {
		case MLGet(MLTNoPermission):
			authorityValue = 1
		case MLGet(MLTReadOnly):
			authorityValue = 2
		case MLGet(MLTReadWrite):
			authorityValue = 3
		}
	})
	authority.Horizontal = true
	authority.SetSelected(MLGet(MLTReadWrite))
	authority.Required = true
	items := []*widget.FormItem{
		{Text: MLGet(MLTInputTitle), Widget: name, HintText: MLGet(MLTInputTitleHint)},
		{Text: MLGet(MLTInputAuthority), Widget: authority, HintText: MLGet(MLTInputAuthorityHint)},
		{Text: MLGet(MLTInputPassword), Widget: password, HintText: MLGet(MLTInputPasswordHint)},
	}

	dialogTitle := MLGet(MLTNewMemo)
	dialogConfirm := MLGet(MLTCreate)
	if isEdit {
		dialogTitle = MLGet(MLTEditsMemo, oldMemoData.Name)
		dialogConfirm = MLGet(MLTSaveEdits)
		name.SetText(oldMemoData.Name)
		password.SetText(oldMemoData.Password)
		authority.SetSelected(authorityMap[oldMemoData.Authority])
	}

	passwordDialog := dialog.NewForm(dialogTitle, dialogConfirm, MLGet(MLTCancel), items, func(b bool) {

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
		defer func() {
			_ = reader.Close()
		}()

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
	fd.Resize(fyne.NewSize(960, 700))
	fd.Show()
}

func CopyMemoEvent(url string) {
	if url == "" {
		dialog.ShowError(fmt.Errorf("复制失败，链接为空"), MainWindow)
		return
	}

	i, ok := define.ShareHas[url]
	if !ok {
		define.ShareId++
		define.ShareHas[url] = define.ShareId
		define.ShareMap[define.ShareId] = url
		i = define.ShareId
	}
	url = fmt.Sprintf("%s/s/%d", define.DoMain, i)

	clipboard := MainApp.Clipboard()
	clipboard.SetContent(url)
	DialogCopySuccess(url)
}

func OpenMemoEvent(memoUrl string) {
	qrImg, _ := utils.GetQRCodeIO(memoUrl)
	reader := bytes.NewReader(qrImg)
	DownloadQr := canvas.NewImageFromReader(reader, "")
	DownloadQr.FillMode = canvas.ImageFillOriginal
	qrDialog := dialog.NewCustom(MLGet(MLTScanQr), MLGet(MLTClose), container.NewCenter(DownloadQr), MainWindow)
	qrDialog.Resize(fyne.NewSize(500, 600))
	qrDialog.Show()
}

func DelMemoEvent() {
	dialog.ShowConfirm(MLGet(MLTDialogTipTitle), MLGet(MLTConfirmDeletion), func(b bool) {
		if b {
			oldMemoData, err := data.GetMemoInfo(NowMemoId)
			if err != nil {
				dialog.ShowError(err, MainWindow)
				return
			}

			err = data.DeleteMemo(NowMemoId)
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
		defer func() {
			_ = file.Close()
		}()

		_, err = file.WriteString(MemoEntry.Text)
		if err != nil {
			logger.Error("写入文件失败：", err)
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
			Event: fmt.Sprintf("另存了备忘录:%s -> %s ", oldMemoData.Name, writer.URI().Path()),
		})
		dialog.ShowInformation(MLGet(MLTDialogTipTitle), MLGet(MLTSaveMemoSuccess, writer.URI().Path()), MainWindow)
	}, MainWindow)
	fd.SetFilter(storage.NewExtensionFileFilter([]string{".txt"}))
	fd.SetFileName(memoData.Name + ".txt")
	fd.SetTitleText("另存为txt")
	fd.Resize(fyne.NewSize(960, 700))
	fd.Show()
}

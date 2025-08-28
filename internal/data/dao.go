package data

import (
	"TFLanHttpDesktop/common/logger"
	"TFLanHttpDesktop/common/utils"
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"time"
)

// GetDownloadData 获取当前下载文件数据
func GetDownloadData() (*DownloadNow, error) {
	result := &DownloadNow{}
	err := DB.Get(DownloadNowTable, DownloadNowTableKey, &result)
	return result, err
}

// GetUploadData 获取当前上传文件数据
func GetUploadData() (*UploadNow, error) {
	result := &UploadNow{}
	err := DB.Get(UploadNowTable, UploadNowTableKey, &result)
	return result, err
}

// todo... 创建当前下载文件数据，不存在创建，存在更新
func SetDownloadData(value *DownloadNow) error {
	return DB.Set(DownloadNowTable, DownloadNowTableKey, &value)
}

// todo... 创建当前上传文件路径，不存在创建，存在更新

// todo... 删除当前下载文件

// todo... 删除当前上传文件

// SetDownloadLog 记录下载日志
func SetDownloadLog(value *DownloadLog) error {
	return DB.Set(DownloadLogTable, utils.AnyToString(time.Now().Unix()), &value)
}

// todo... 记录上传日志

// todo... 查看下载日志
func GetDownloadLog() ([]*DownloadLog, error) {
	limit := 1000
	result := make([]*DownloadLog, 0)
	db := DB.GetDB()
	defer db.Close()
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(DownloadLogTable))
		if bucket == nil {
			return fmt.Errorf("bucket %s 不存在", DownloadLogTable)
		}

		cursor := bucket.Cursor()
		count := 0

		// 1. 移动到最后一个key
		k, v := cursor.Last()
		for k != nil && count < limit {
			item := &DownloadLog{}
			err := json.Unmarshal(v, &item)
			if err != nil {
				log.Fatal(err)
			} else {
				result = append(result, item)
			}
			// 3. 向前移动游标
			k, v = cursor.Prev()
			count++
		}
		return nil
	})
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return result, nil
}

// todo... 查看上传日志

// todo... 记录操作日志

// todo... 查看操作日志

// todo... 创建备忘录

// todo... 修改备忘录内容

// todo... 查看备忘录信息

// todo... 查看备忘录内容

// todo... 删除备忘录

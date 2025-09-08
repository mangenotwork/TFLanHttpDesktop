package data

import (
	"TFLanHttpDesktop/common/logger"
	"TFLanHttpDesktop/common/utils"
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"sort"
	"strings"
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

// SetDownloadData 创建当前下载文件数据，不存在创建，存在更新
func SetDownloadData(value *DownloadNow) error {
	return DB.Set(DownloadNowTable, DownloadNowTableKey, &value)
}

// SetUploadData 创建当前上传文件路径，不存在创建，存在更新
func SetUploadData(value *UploadNow) error {
	return DB.Set(UploadNowTable, UploadNowTableKey, &value)
}

// SetDownloadLog 记录下载日志
func SetDownloadLog(value *DownloadLog) error {
	return DB.Set(DownloadLogTable, utils.AnyToString(time.Now().Unix()), &value)
}

// SetUploadLog 记录上传日志
func SetUploadLog(value *UploadLog) error {
	return DB.Set(UploadLogTable, utils.AnyToString(time.Now().Unix()), &value)
}

// GetDownloadLog 查看下载日志
func GetDownloadLog() ([]*DownloadLog, error) {
	limit := 1000
	result := make([]*DownloadLog, 0)
	db := DB.GetDB()
	defer func() {
		_ = db.Close()
	}()
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

// GetUploadLog 查看上传日志
func GetUploadLog() ([]*UploadLog, error) {
	limit := 1000
	result := make([]*UploadLog, 0)
	db := DB.GetDB()
	defer func() {
		_ = db.Close()
	}()
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(UploadLogTable))
		if bucket == nil {
			return fmt.Errorf("bucket %s 不存在", UploadLogTable)
		}

		cursor := bucket.Cursor()
		count := 0

		// 1. 移动到最后一个key
		k, v := cursor.Last()
		for k != nil && count < limit {
			item := &UploadLog{}
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

// GetMemoList 获取备忘录列表
func GetMemoList() ([]*Memo, error) {
	result := make([]*Memo, 0)
	err := DB.GetAll(MemoTable, func(k, v []byte) {
		item := &Memo{}
		err := json.Unmarshal(v, &item)
		if err != nil {
			logger.Error(err)
		} else {
			result = append(result, item)
		}
	})
	return result, err
}

// NewMemo 创建备忘录
func NewMemo(name string, authority int, password string) (*Memo, error) {
	id := utils.IDMd5()
	now := time.Now()
	isPassword := 0
	if password != "" {
		isPassword = 1
	}
	memo := &Memo{
		Id:         id,
		Name:       name,
		CreateTime: now,
		LastTime:   now,
		Authority:  authority,
		IsPassword: isPassword,
		Password:   password,
	}
	err := DB.Set(MemoTable, id, memo)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	err = DB.Set(MemoContentTable, id, MemoContent(""))
	return memo, err
}

func GetMemoContent(id string) (MemoContent, error) {
	var content MemoContent = ""
	err := DB.Get(MemoContentTable, id, &content)
	return content, err
}

// SetMemoContent 修改备忘录内容
func SetMemoContent(id string, content string) (MemoContent, error) {
	err := DB.Set(MemoContentTable, id, content)
	go func() {
		_ = GetMemoFenCiList(id, content)
	}()

	//now := time.Now()
	//t := now.Sub(MemoEntryTime).Seconds()
	//if t > 0.02 {
	//	go func() {
	//		_ = GetMemoFenCiList(id, content)
	//	}()
	//}
	//logger.Debug("写入时差 : ", t)
	//MemoEntryTime = now

	return MemoContent(content), err
}

// GetMemoInfo 查看备忘录信息
func GetMemoInfo(id string) (*Memo, error) {
	info := &Memo{}
	err := DB.Get(MemoTable, id, &info)
	return info, err
}

// SetMemoInfo 查看备忘录信息
func SetMemoInfo(id string, name string, authority int, password string) (*Memo, error) {
	now := time.Now()
	nowMemo, err := GetMemoInfo(id)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	nowMemo.LastTime = now

	if name != "" {
		nowMemo.Name = name
	}

	if authority != 0 {
		nowMemo.Authority = authority
	}

	nowMemo.Password = password
	if nowMemo.Password == "" {
		nowMemo.IsPassword = 0
	} else {
		nowMemo.IsPassword = 1
	}
	err = DB.Set(MemoTable, id, nowMemo)
	return nowMemo, err
}

// DeleteMemo 删除备忘录
func DeleteMemo(id string) error {
	logger.Debug("DeleteMemo = ", id)
	err := DB.Delete(MemoTable, id)
	if err != nil {
		logger.Error(err)
		return err
	}
	err = DB.Delete(MemoContentTable, id)
	if err != nil {
		logger.Error(err)
		return err
	}

	err = DelMemoCiList(id)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// SetOperationLog 记录操作日志
func SetOperationLog(value *OperationLog) error {
	return DB.Set(OperationLogTable, utils.AnyToString(time.Now().Unix()), &value)
}

// GetOperationLog 查看操作日志
func GetOperationLog() ([]*OperationLog, error) {
	limit := 1000
	result := make([]*OperationLog, 0)
	db := DB.GetDB()
	defer func() {
		_ = db.Close()
	}()
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(OperationLogTable))
		if bucket == nil {
			return fmt.Errorf("bucket %s 不存在", OperationLogTable)
		}

		cursor := bucket.Cursor()
		count := 0

		// 1. 移动到最后一个key
		k, v := cursor.Last()
		for k != nil && count < limit {
			item := &OperationLog{}
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

func GetMemoCiList(id string) ([]*MemoCiList, error) {
	result := make([]*MemoCiList, 0)
	err := FcDB.Get(MemoCiListTable, id, &result)
	return result, err
}

func SetMemoCiList(id string, list []*MemoCiList) error {
	err := FcDB.Set(MemoCiListTable, id, &list)
	return err
}

func DelMemoCiList(id string) error {
	ciList, err := GetMemoCiList(id)
	if err != nil {
		logger.Error(err)
		return err
	}
	for _, ci := range ciList {
		_ = DelCiList(ci.Ci, id)
	}
	return FcDB.Delete(MemoCiListTable, id)
}

func MatchCi(ci string) []string {
	result := make([]string, 0)
	list, err := CiDB.AllKey(CiListTable)
	if err != nil {
		logger.Error(err)
		return []string{}
	}
	for _, v := range list {
		if strings.Contains(v, ci) {
			result = append(result, v)
		}
	}
	return result
}

func GetCiList(ci string) ([]*CiList, error) {
	result := make([]*CiList, 0)
	err := CiDB.Get(CiListTable, ci, &result)
	sort.Slice(result, func(i, j int) bool {
		if result[i].WordFrequency > result[j].WordFrequency {
			return true
		}
		return false
	})
	return result, err
}

func SetCiList(ci string, list []*CiList) error {
	err := CiDB.Set(CiListTable, ci, &list)
	return err
}

func DelCiList(ci string, memoId string) error {
	list, _ := GetCiList(ci)

	for i := len(list) - 1; i >= 0; i-- {
		if list[i].MemoId == memoId {
			list = append(list[:i], list[i+1:]...)
		}
	}

	return SetCiList(ci, list)
}

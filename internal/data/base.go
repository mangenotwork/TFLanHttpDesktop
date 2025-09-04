package data

import (
	"TFLanHttpDesktop/common/logger"
	"TFLanHttpDesktop/common/utils"
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"os"
	"strings"
	"time"
)

const (
	DownloadNowTable    = "DownloadNowTable" // DownloadNow
	DownloadNowTableKey = "DownloadNowTableKey"
	DownloadLogTable    = "DownloadLogTable" // DownloadLog
	UploadNowTable      = "UploadNowTable"   // UploadNow
	UploadNowTableKey   = "UploadNowTableKey"
	UploadLogTable      = "UploadLogTable"    // UploadLog
	MemoTable           = "MemoTable"         // Memo
	MemoContentTable    = "MemoContentTable"  // MemoContent
	OperationLogTable   = "OperationLogTable" // OperationLog
	MemoCiListTable     = "MemoCiListTable"   // MemoCiList
	CiListTable         = "CiListTable"       // CiList
)

var DB *LocalDB
var FcDB *LocalDB
var CiDB *LocalDB

var Tables = []string{DownloadNowTable, DownloadLogTable, UploadNowTable, UploadLogTable,
	MemoTable, MemoContentTable, OperationLogTable}
var FcTables = []string{MemoCiListTable}
var CiTables = []string{CiListTable}
var ISNULL = fmt.Errorf("ISNULL")
var TableNotFound = fmt.Errorf("table notfound")

type LocalDB struct {
	Path   string
	Tables []string
	Conn   *bolt.DB
}

func getDirName(filePath string) string {
	return filePath[0:getLastIndex(filePath, `/`)]
}

func getLastIndex(str, ch string) int {
	return len(str) - len(strings.TrimRight(str, ch))
}

func checkDBFile(dbFilePath string) {

	if _, err := os.Stat(dbFilePath); os.IsNotExist(err) {

		dir := getDirName(dbFilePath)

		if dir != "" {
			if err = os.MkdirAll(dir, 0755); err != nil {
				log.Panic(err)
			}
		}

		f, fErr := os.Create(dbFilePath)
		if fErr != nil {
			log.Panic(fErr)
		}

		defer func() {
			_ = f.Close()
		}()

	}

}

func (ldb *LocalDB) Init() {

	checkDBFile(ldb.Path)

	db, err := bolt.Open(ldb.Path, 0600, nil)
	if err != nil {
		logger.Panic(err)
	}

	defer func() {
		_ = db.Close()
	}()

	for _, table := range ldb.Tables {

		err = db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(table))

			if b == nil {
				_, err = tx.CreateBucket([]byte(table))
				if err != nil {
					logger.Panic(err)
				}
			}

			return nil
		})

		if err != nil {
			logger.Panic(err)
		}

	}
}

func NewLocalDB(tables []string, path string) *LocalDB {
	return &LocalDB{
		Path:   path,
		Tables: tables,
	}
}

var MemoEntryTime time.Time

func InitDB(dbPath, fcPath, ciPath string) {
	DB = NewLocalDB(Tables, dbPath)
	DB.Init()
	FcDB = NewLocalDB(FcTables, fcPath)
	FcDB.Init()
	CiDB = NewLocalDB(CiTables, ciPath)
	CiDB.Init()
	// 初始化分词
	Seg.SkipLog = true
	err := Seg.LoadDict()
	if err != nil {
		log.Panic("初始化分词失败")
	}
	MemoEntryTime = time.Now()
}

func (ldb *LocalDB) Open() {
	var err error
	maxRetries := 4
	for i := 0; i < maxRetries; i++ {
		ldb.Conn, err = bolt.Open(ldb.Path, 0600, &bolt.Options{
			Timeout: 5 * time.Second, // 设置超时时间，避免永久阻塞
		})
		if err == nil {
			break
		}
		time.Sleep(40 * time.Millisecond)
	}
}

func (ldb *LocalDB) Close() {
	if ldb.Conn != nil {
		_ = ldb.Conn.Close()
	}
}

func (ldb *LocalDB) GetDB() *bolt.DB {
	ldb.Open()
	return ldb.Conn
}

func (ldb *LocalDB) ClearTable(table string) error {
	ldb.Open()
	defer func() {
		_ = ldb.Conn.Close()
	}()
	return ldb.Conn.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(table))
	})
}

func (ldb *LocalDB) Stats(table string) (bolt.BucketStats, error) {

	var stats bolt.BucketStats

	ldb.Open()

	defer func() {
		if ldb.Conn != nil {
			_ = ldb.Conn.Close()
		}
	}()

	err := ldb.Conn.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(table))
		if b == nil {
			err := ldb.ClearTable(table)
			if err != nil {
				return err
			}
		}

		stats = b.Stats()

		return nil
	})

	return stats, err
}

func (ldb *LocalDB) Get(table, key string, data interface{}) error {

	ldb.Open()
	defer func() {
		if ldb.Conn != nil {
			_ = ldb.Conn.Close()
		}
	}()

	return ldb.Conn.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(table))
		if b == nil {
			return TableNotFound
		}

		bt := b.Get([]byte(key))
		if len(bt) < 1 {
			return ISNULL
		}

		err := json.Unmarshal(bt, data)
		if err != nil {
			return err
		}

		return nil
	})
}

func (ldb *LocalDB) Set(table, key string, data interface{}) error {

	value, err := utils.AnyToJsonB(data)
	if err != nil {
		return err
	}

	ldb.Open()

	defer func() {
		if ldb.Conn != nil {
			_ = ldb.Conn.Close()
		}
	}()

	if ldb.Conn == nil {
		logger.Error("未取到本地数据连接")
		return fmt.Errorf("未取到本地数据连接")
	}

	return ldb.Conn.Update(func(tx *bolt.Tx) error {

	R:
		b := tx.Bucket([]byte(table))
		if b == nil {
			_, err = tx.CreateBucket([]byte(table))
			if err != nil {
				return err
			}

			goto R
		}

		err = b.Put([]byte(key), value)
		if err != nil {
			return err
		}

		return nil
	})
}

func (ldb *LocalDB) Delete(table, key string) error {
	ldb.Open()

	defer func() {
		if ldb.Conn != nil {
			_ = ldb.Conn.Close()
		}
	}()

	return ldb.Conn.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(table))
		if b == nil {
			return fmt.Errorf("未获取到表")
		}
		if err := b.Delete([]byte(key)); err != nil {
			return err
		}
		return nil
	})
}

func (ldb *LocalDB) AllKey(table string) ([]string, error) {
	keys := make([]string, 0)

	ldb.Open()

	defer func() {
		if ldb.Conn != nil {
			_ = ldb.Conn.Close()
		}
	}()

	err := ldb.Conn.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(table))
		if b == nil {
			return TableNotFound
		}

		c := b.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			keys = append(keys, string(k))
		}

		return nil
	})
	return keys, err
}

func (ldb *LocalDB) GetAll(table string, f func(k, v []byte)) error {
	ldb.Open()

	defer func() {
		if ldb.Conn != nil {
			_ = ldb.Conn.Close()
		}
	}()

	err := ldb.Conn.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(table))
		if b == nil {
			return TableNotFound
		}

		return b.ForEach(func(k, v []byte) error {
			f(k, v)
			return nil
		})

	})
	return err
}

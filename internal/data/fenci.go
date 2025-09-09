package data

import (
	"TFLanHttpDesktop/common/logger"
	"errors"
	"github.com/go-ego/gse"
	"strings"
)

var (
	Seg gse.Segmenter
)

type Term struct {
	Text  string // 词
	Freq  float64
	Start int // 开始段
	End   int // 结束段
	Pos   string
}

// TermExtract 提取索引词
// 除了标点符号，助词，语气词，形容词，叹词, 副词 其他都被分出来
func TermExtract(str string) []*Term {
	segments := Seg.Segment([]byte(str))
	termList := make([]*Term, 0)
	for _, v := range segments {
		t := v.Token()
		p := t.Pos()
		txt := t.Text()
		end := v.End()
		start := v.Start()
		//logger.Info("txt = ", txt, p)

		if p == "w" || p == "u" || p == "uj" || p == "y" || p == "a" || p == "e" || p == "d" {
			continue
		}

		if p == "x" && !ContainsEnglishAndNumber(txt) {
			continue
		}

		termList = append(termList, &Term{
			Text:  txt,
			Freq:  t.Freq(),
			End:   end,
			Start: start,
			Pos:   p,
		})
	}
	return termList
}

func ContainsEnglishAndNumber(str string) bool {
	dictionary := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	for _, v := range str {
		if strings.Contains(dictionary, string(v)) {
			return true
		}
	}
	return false
}

// MemoFenCiList 词:词频
type MemoFenCiList map[string]int

// ANotB 找出当前map(A)中存在、但参数map(B)中不存在的词，返回新的差集map
func (a MemoFenCiList) ANotB(b MemoFenCiList) MemoFenCiList {
	notMap := make(MemoFenCiList)
	for word, count := range a {
		if _, exists := b[word]; !exists {
			notMap[word] = count
		}
	}
	return notMap
}

func GetMemoFenCiList(memoId, content string) error {

	oldFenCi, err := GetMemoCiList(memoId)
	if err != nil && !errors.Is(ISNULL, err) {
		logger.Error("获取备忘录分词数据失败: ", err)
		return err
	}

	oldFenCiMap := make(MemoFenCiList)
	for _, v := range oldFenCi {
		oldFenCiMap[v.Ci] = v.WordFrequency
	}
	//logger.Debug("oldFenCiMap = ", oldFenCiMap)

	memoFc := make(MemoFenCiList)
	fcList := TermExtract(content)
	for _, v := range fcList {
		if _, ok := memoFc[v.Text]; !ok {
			memoFc[v.Text] = 1
		} else {
			memoFc[v.Text]++
		}
	}
	//logger.Debug("memoFc = ", memoFc)

	// 删除
	delFcList := oldFenCiMap.ANotB(memoFc)
	//logger.Debug("delFcList = ", delFcList)
	for k, _ := range delFcList {
		fc1, err := GetCiList(k)
		if err != nil && !errors.Is(ISNULL, err) {
			logger.Error("获取词失败:", err)
			continue
		}
		for i := len(fc1) - 1; i >= 0; i-- {
			if fc1[i].MemoId == memoId {
				fc1 = append(fc1[:i], fc1[i+1:]...)
			}
		}
		err = SetCiList(k, fc1)
		if err != nil {
			logger.Error("保存删除后的词失败:", err)
		}
	}

	// 新增
	addFcList := memoFc.ANotB(oldFenCiMap)
	//logger.Debug("addFcList = ", addFcList)
	for k, v := range addFcList {
		fc1, err := GetCiList(k)
		if err != nil && !errors.Is(ISNULL, err) {
			logger.Error("获取词失败:", err)
			continue
		}
		fc1 = append(fc1, &CiList{
			MemoId:        memoId,
			WordFrequency: v,
		})
		err = SetCiList(k, fc1)
		if err != nil {
			logger.Error("保存删除后的词失败:", err)
		}
	}

	// 记录
	newFenCi := make([]*MemoCiList, 0)
	for k, v := range memoFc {
		newFenCi = append(newFenCi, &MemoCiList{
			Ci:            k,
			WordFrequency: v,
		})
	}
	err = SetMemoCiList(memoId, newFenCi)
	if err != nil {
		logger.Error("记录新的分词失败")
		return err
	}
	return nil
}

package ioms

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

type CheckFile struct {
	n        int
	FileName string
	Total    int
	e        bool
	lock     sync.Locker
	FullPath string
}
type CheckItem struct {
	BankAcct    string
	ExchAcct    string
	BankSid     string
	ExchSid     string
	TrDate      string
	TrTime      string
	TrType      string
	TrAmt       string
	Currency    string
	CashOrRemit string
	Initiator   string
	AdjustFlag  string
	IdType      string
	IdNum       string
	CustName    string
	Reserved    string
	TransFlag   string // 交易标记 1 成功， 2 失败
}

func NewCheckFile(path string, bankid int, batchno int, date string, total int) *CheckFile {

	fileName := fmt.Sprintf("%d_%s_CHECK_%d", bankid, date, batchno)
	filePath := fmt.Sprintf("%s%s", path, fileName)

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return nil
	}
	defer file.Close()

	f := &CheckFile{
		FullPath: filePath,
		n:        0,
		e:        false,
		Total:    total,
		FileName: fileName,
		lock:     &sync.Mutex{},
	}
	return f
}
func (f *CheckFile) CheckDone() bool {
	f.lock.Lock()
	defer f.lock.Unlock()

	return f.Total <= f.n

}
func (f *CheckFile) ReduceCheckCount(n int) int {
	f.lock.Lock()
	defer f.lock.Unlock()

	f.Total = f.Total - n
	return f.Total
}
func (f *CheckFile) Append(item *CheckItem) (int, error) {
	f.lock.Lock()
	defer f.lock.Unlock()
	if item == nil {
		// 表示处理中的
		f.Total = f.Total - 1
		return f.n, nil
	}
	str := fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s\n",
		item.BankAcct, item.ExchAcct, item.BankSid, item.ExchSid, item.TrDate,
		item.TrTime, item.TrType, item.TrAmt, item.Currency, item.CashOrRemit,
		item.Initiator, item.AdjustFlag, item.IdType, item.IdNum, item.Currency,
		item.Reserved, item.TransFlag,
	)
	err := ioutil.WriteFile(f.FullPath, []byte(str), 0664)
	if err != nil {
		f.e = true
	}
	if f.e {
		return 0, fmt.Errorf("some go routine write file error, pls recheck manual")
	}
	f.n = f.n + 1

	return f.n, nil
}

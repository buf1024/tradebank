package util

import (
	"fmt"
)

const (
	E_SUCCESS      = 99999
	E_BANK_ERR     = 15973
	E_HALF_SUCCESS = 44444
)

type TradeError struct {
	Code int64
}

var err map[int64]string

func (e TradeError) Error() string {
	if msg, ok := err[e.Code]; ok {
		return fmt.Sprintf("[ERR=%d, EMSG=%s]", e.Code, msg)
	}
	return fmt.Sprintf("[Not Found]")
}

func NewError(code int64) TradeError {
	e := TradeError{}
	e.Code = code
	return e
}
func GetErrMsg(code int64) string {
	e := NewError(code)
	return e.Error()
}

func init() {
	err = make(map[int64]string)

	err[E_SUCCESS] = "处理成功"

}

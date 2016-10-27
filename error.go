package tradebank

import (
	"fmt"
)

const (
	E_SUCCESS = 99999
)

type TradeError struct {
	Code    int64
	Message string
}

var err map[int64]string

func (t *TradeError) GetError(code int64) (string, error) {
	msg, ok := err[code]
	if ok {
		return msg, nil
	}
	return "", fmt.Errorf("Not Found")
}

func init() {
	err[E_SUCCESS] = "处理成功"

}

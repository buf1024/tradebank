package yoyi

import (
	"tradebank/ioms"
)

type yoyi struct {
}

func (b *yoyi) Name() string {
	return "YOYI"
}

func (b *yoyi) ID() int64 {
	return 26
}

func (b *yoyi) LoadConfig(path string) error {
	return nil
}

func init() {
	b := &yoyi{}
	RegisterBank(b)
}

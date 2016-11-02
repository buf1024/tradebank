package bank

import (
	"fmt"
)

type MyBank interface {
	Name() string
	ID() int64
	Init(path string) error
}

var mybank map[int64]string

func Bank(id int64) (MyBank, error) {
	if bank, ok := mybank[b.ID()]; ok {
		return bank, nil
	}
	return nil, fmt.Errorf("bank %d not register", id)
}

func Register(b MyBank) {
	if bank, ok := mybank[b.ID()]; !ok {
		mybank[b.ID()] = b
	}
}

func init() {
	mybank = make(map[int64]bank.MyBank)
}

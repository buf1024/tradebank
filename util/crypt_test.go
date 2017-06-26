package util

import (
	"fmt"
	"testing"
)

func TestCrypt(t *testing.T) {
	c := Crypt{}
	encStr := "cmnsf+KOk/kzY5m7nIBCzA=="
	decStr, err := c.DBDecrypt(encStr)
	if err != nil {
		fmt.Printf("ERR=%s\n", err)
		t.FailNow()
	}
	fmt.Printf("decrypt:= %s\n", decStr)

	newEncStr, err := c.DBEncrypt(decStr)
	if err != nil {
		fmt.Printf("ERR=%s\n", err)
		t.FailNow()
	}
	fmt.Printf("encrypt:= %s\n", newEncStr)
	if encStr != newEncStr {
		t.FailNow()
	}
}

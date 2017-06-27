package util

import (
	"encoding/base64"
	"fmt"
	"testing"
)

func TestDBCrypt(t *testing.T) {
	encStr := "cmnsf+KOk/kzY5m7nIBCzA=="
	decStr, err := DBDecrypt(encStr)
	if err != nil {
		fmt.Printf("ERR=%s\n", err)
		t.FailNow()
	}
	fmt.Printf("decrypt:= %s\n", decStr)

	newEncStr, err := DBEncrypt(decStr)
	if err != nil {
		fmt.Printf("ERR=%s\n", err)
		t.FailNow()
	}
	fmt.Printf("encrypt:= %s\n", newEncStr)
	if encStr != newEncStr {
		t.FailNow()
	}
}

func TestPacketCrypt(t *testing.T) {
	key := "Q09NLkNPTkdYSU5H"
	rawStr := "12345678---===---"
	encStr, err := PacketEncrypt(key, []byte(rawStr))
	if err != nil {
		fmt.Printf("ERR=%s\n", err)
		t.FailNow()
	}
	fmt.Printf("encrypt base64:= %s\n", base64.StdEncoding.EncodeToString(encStr))

	decByte, err := PacketDecrypt(key, []byte(encStr))
	if err != nil {
		fmt.Printf("ERR=%s\n", err)
		t.FailNow()
	}

	// X/M0ZQcfCYYLgp4wqf7xWg==
	// X/M0ZQcfCYaXLrSRzl4Uuw==
	fmt.Printf("decrypt:= %s\n", string(decByte))
}

package main

import (
	"bytes"
	"crypto/aes"
	"encoding/hex"
	"strings"
)

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5Unpadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
func NoCardPayEncrypt(src string, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize()
	srcPadding := PKCS5Padding([]byte(src), blockSize)
	dstEnc := make([]byte, len(srcPadding))
	for bs, be := 0, blockSize; bs < len(srcPadding); bs, be = bs+blockSize, be+blockSize {
		block.Encrypt(dstEnc[bs:be], srcPadding[bs:be])
	}
	dstStr := hex.EncodeToString(dstEnc)
	return strings.ToUpper(dstStr), nil
}
func NoCardPayDecrypt(src string, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	data, err := hex.DecodeString(src)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	decBytes := make([]byte, len(data))
	srcBytes := []byte(data)
	for bs, be := 0, blockSize; bs < len(data); bs, be = bs+blockSize, be+blockSize {
		block.Decrypt(decBytes[bs:be], srcBytes[bs:be])
	}
	decStr := string(PKCS5Unpadding(decBytes))
	return decStr, nil
}

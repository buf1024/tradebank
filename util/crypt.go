package util

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
	"fmt"
	"math"
)

/*
type Caesar struct {
	UC_ENCRYPT_CHARS []byte
	LC_ENCRYPT_CHARS []byte
	UC_DECRYPT_CHARS []byte
	LC_DECRYPT_CHARS []byte
}

func (c *Caesar) Init() {
	UC_DECRYPT_CHARS = make([]byte, 26)
	LC_DECRYPT_CHARS = make([]byte, 26)

	UC_ENCRYPT_CHARS = []byte{'M', 'D', 'X', 'U', 'P', 'I', 'B', 'E', 'J', 'C', 'T', 'N',
		'K', 'O', 'G', 'W', 'R', 'S', 'F', 'Y', 'V', 'L', 'Z', 'Q', 'A', 'H'}

	LC_ENCRYPT_CHARS = []byte{'m', 'd', 'x', 'u', 'p', 'i', 'b', 'e', 'j', 'c', 't', 'n',
		'k', 'o', 'g', 'w', 'r', 's', 'f', 'y', 'v', 'l', 'z', 'q', 'a', 'h'}
	for i := 0; i < 26; i++ {
		b := UC_ENCRYPT_CHARS[i]
		UC_DECRYPT_CHARS[b-'A'] = (byte)('A' + i)
		b = LC_ENCRYPT_CHARS[i]
		LC_DECRYPT_CHARS[b-'a'] = (byte)('a' + i)
	}
}
func (c *Caesar) Encrypt(b byte) byte {
	if b >= 'A' && b <= 'Z' {
		return UC_ENCRYPT_CHARS[b-'A']
	} else if b >= 'a' && b <= 'z' {

		return LC_ENCRYPT_CHARS[b-'a']
	}

	return b
}
func (c *Caesar) Decrypt(b byte) byte {
	if b >= 'A' && b <= 'Z' {
		return UC_DECRYPT_CHARS[b-'A']
	} else if b >= 'a' && b <= 'z' {
		return LC_DECRYPT_CHARS[b-'a']
	}
	return b
}
func (c *Caesar) EncryptBytes(b []byte) []byte {
	dst := make([]byte, len(b))
	for i := 0; i < len(b); i++ {
		dst[i] = Encrypt(b[i])
	}
	return dst
}
func (c *Caesar) DecryptBytes(b []byte) []byte {
	dst := make([]byte, len(b))
	for i := 0; i < len(b); i++ {
		dst[i] = Decrypt(b[i])
	}
	return dst
}
*/

var des3DefKey []byte = []byte("ZGFyayBtYXR0ZXIgZXNjYXBl")
var des3DefIov []byte = []byte("bmV2ZXI=")

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
func getDBKey() string {
	s := math.Sqrt(17.0)
	key := fmt.Sprintf("%.22f", s)
	return key[:8]
}
func getPacketKey(key string) ([]byte, []byte, error) {
	keyByte := make([]byte, 24)
	if len(key) > 0 {
		len := len(key)
		if len < 8 || len > 24 {
			return nil, nil, fmt.Errorf("key's len must >=8 && <= 24")
		}
		copy(keyByte, []byte(key))
		rst := 24 - len
		for i := 0; i < rst; i++ {
			keyByte[len+i] = byte(0)
		}
	}
	return keyByte, des3DefIov, nil

}
func noStdPadding(ciphertext []byte, blockSize int) []byte {
	bodylen := len(ciphertext)
	append := bodylen % blockSize
	if append > 0 {
		append = 8 - append
		bodylen = bodylen + append
	}
	padding := make([]byte, bodylen)
	copy(padding, ciphertext)
	for i := 0; i < append; i++ {
		padding[bodylen-append+i] = byte(0)
	}
	return padding
}

func noStdUnpadding(origData []byte) []byte {
	length := len(origData)
	unpadding := 0
	for i := length - 1; i >= length-8; i-- {
		if origData[i] == byte(0) {
			unpadding++
		} else {
			break
		}
	}
	return origData[:(length - unpadding)]
}

func PacketEncrypt(key string, buf []byte) (dst []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			dst = nil
			err = fmt.Errorf("%s", e)
		}
	}()
	keyBytes, iovBytes, err := getPacketKey(key)
	if err != nil {
		return nil, err
	}
	block, err := des.NewTripleDESCipher(keyBytes)
	if err != nil {
		return nil, err
	}
	bufBase64 := base64.StdEncoding.EncodeToString(buf)
	padBytes := noStdPadding([]byte(bufBase64), block.BlockSize())
	mode := cipher.NewCBCEncrypter(block, iovBytes)
	dst = make([]byte, len(padBytes))
	mode.CryptBlocks(dst, padBytes)

	return dst, nil
}

func PacketDecrypt(key string, buf []byte) (decBytes []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			decBytes = nil
			err = fmt.Errorf("%s", e)
		}
	}()
	keyBytes, iovBytes, err := getPacketKey(key)
	if err != nil {
		return nil, err
	}
	block, err := des.NewTripleDESCipher(keyBytes)
	if err != nil {
		return nil, err
	}
	mode := cipher.NewCBCDecrypter(block, iovBytes)
	dst := make([]byte, len(buf))
	mode.CryptBlocks(dst, buf)
	unpadBytes := noStdUnpadding(dst)
	decBytes, err = base64.StdEncoding.DecodeString(string(unpadBytes))
	if err != nil {
		return nil, err
	}

	return decBytes, nil
}

func PswEncrypt(pass string, signacct string) (string, error) {
	return "buf", nil
}

func PswDecrypt(encstr string, signacct string) (string, error) {
	return "buf", nil
}

func DBDecrypt(buf string) (string, error) {
	orig, err := base64.StdEncoding.DecodeString(buf)
	if err != nil {
		return "", nil
	}
	key := getDBKey()
	block, err := des.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	b := new(bytes.Buffer)
	for i := 0; i < len(key); i++ {
		b.WriteByte(byte(0))
	}
	mode := cipher.NewCBCDecrypter(block, b.Bytes())
	dst := make([]byte, len(orig))
	mode.CryptBlocks(dst, orig)
	dst = PKCS5Unpadding(dst)
	return string(dst), nil
}

func DBEncrypt(buf string) (string, error) {
	key := getDBKey()
	block, err := des.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	b := new(bytes.Buffer)
	for i := 0; i < len(key); i++ {
		b.WriteByte(byte(0))
	}
	mode := cipher.NewCBCEncrypter(block, b.Bytes())
	bufPad := PKCS5Padding([]byte(buf), mode.BlockSize())
	dst := make([]byte, len(bufPad))
	mode.CryptBlocks(dst, bufPad)
	return base64.StdEncoding.EncodeToString(dst), nil
}

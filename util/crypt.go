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
	c.UC_DECRYPT_CHARS = make([]byte, 26)
	c.LC_DECRYPT_CHARS = make([]byte, 26)

	c.UC_ENCRYPT_CHARS = []byte{'M', 'D', 'X', 'U', 'P', 'I', 'B', 'E', 'J', 'C', 'T', 'N',
		'K', 'O', 'G', 'W', 'R', 'S', 'F', 'Y', 'V', 'L', 'Z', 'Q', 'A', 'H'}

	c.LC_ENCRYPT_CHARS = []byte{'m', 'd', 'x', 'u', 'p', 'i', 'b', 'e', 'j', 'c', 't', 'n',
		'k', 'o', 'g', 'w', 'r', 's', 'f', 'y', 'v', 'l', 'z', 'q', 'a', 'h'}
	for i := 0; i < 26; i++ {
		b := c.UC_ENCRYPT_CHARS[i]
		c.UC_DECRYPT_CHARS[b-'A'] = (byte)('A' + i)
		b = c.LC_ENCRYPT_CHARS[i]
		c.LC_DECRYPT_CHARS[b-'a'] = (byte)('a' + i)
	}
}
func (c *Caesar) Encrypt(b byte) byte {
	if b >= 'A' && b <= 'Z' {
		return c.UC_ENCRYPT_CHARS[b-'A']
	} else if b >= 'a' && b <= 'z' {

		return c.LC_ENCRYPT_CHARS[b-'a']
	}

	return b
}
func (c *Caesar) Decrypt(b byte) byte {
	if b >= 'A' && b <= 'Z' {
		return c.UC_DECRYPT_CHARS[b-'A']
	} else if b >= 'a' && b <= 'z' {
		return c.LC_DECRYPT_CHARS[b-'a']
	}
	return b
}
func (c *Caesar) EncryptBytes(b []byte) []byte {
	dst := make([]byte, len(b))
	for i := 0; i < len(b); i++ {
		dst[i] = c.Encrypt(b[i])
	}
	return dst
}
func (c *Caesar) DecryptBytes(b []byte) []byte {
	dst := make([]byte, len(b))
	for i := 0; i < len(b); i++ {
		dst[i] = c.Decrypt(b[i])
	}
	return dst
}
*/
type Crypt struct {
	//cae *Caesar
}

func NewCrypt() *Crypt {
	c := &Crypt{}

	return c
}

func (c *Crypt) getDBKey() string {
	s := math.Sqrt(17.0)
	key := fmt.Sprintf("%.22f", s)
	return key[:8]
}

func (c *Crypt) PacketEncrypt(buf []byte) ([]byte, error) {
	return buf, nil
}

func (c *Crypt) PacketDecrypt(buf []byte) ([]byte, error) {
	return buf, nil
}

func (c *Crypt) PswEncrypt(pass string, signacct string) (string, error) {
	return "buf", nil
}

func (c *Crypt) PswDecrypt(encstr string, signacct string) (string, error) {
	return "buf", nil
}

func (c *Crypt) DBDecrypt(buf string) (string, error) {
	orig, err := base64.StdEncoding.DecodeString(buf)
	if err != nil {
		return "", nil
	}
	key := c.getDBKey()
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
	return string(dst), nil
}

func (c *Crypt) DBEncrypt(buf string) (string, error) {
	key := c.getDBKey()
	block, err := des.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	b := new(bytes.Buffer)
	for i := 0; i < len(key); i++ {
		b.WriteByte(byte(0))
	}
	mode := cipher.NewCBCEncrypter(block, b.Bytes())
	dst := make([]byte, len(buf))
	mode.CryptBlocks(dst, []byte(buf))
	return base64.StdEncoding.EncodeToString(dst), nil
}

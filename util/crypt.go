package util

import (
	"fmt"
	"math"
)

type Crypt struct {
}

func NewCrypt() *Crypt {
	c := &Crypt{}

	return c
}

func (c *Crypt) getDBKey() string {
	s := math.Sqrt(17.0)
	return fmt.Sprintf("%.22f", s)
}

func (c *Crypt) PacketEncrypt(buf []byte) ([]byte, error) {
	return buf, nil
}

func (c *Crypt) PacketDecrypt(buf []byte) ([]byte, error) {
	return buf, nil
}

func (c *Crypt) PswEncrypt(buf []byte) ([]byte, error) {
	return buf, nil
}

func (c *Crypt) PswDecrypt(buf []byte) ([]byte, error) {
	return buf, nil
}

func (c *Crypt) DBDecrypt(buf string) (string, error) {
	return buf, nil
}

func (c *Crypt) DBEncrypt(buf string) (string, error) {
	return buf, nil
}

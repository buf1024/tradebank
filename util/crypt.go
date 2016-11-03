package util

import ()

type Crypt struct {
}

func NewCrypt() *Crypt {
	c := &Crypt{}

	return c
}

func (c *Crypt) PacketEncrypt(buf []byte) ([]byte, error) {
	return buf, nil
}

func (c *Crypt) PacketDecrypt(buf []byte) ([]byte, error) {
	return buf, nil
}

func (c *Crypt) PasswordEncrypt(buf []byte) ([]byte, error) {
	return buf, nil
}

func (c *Crypt) PasswordDecrypt(buf []byte) ([]byte, error) {
	return buf, nil
}

func (c *Crypt) DBStringEncrypt(buf []byte) ([]byte, error) {
	return buf, nil
}

func (c *Crypt) DBStringDecrypt(buf []byte) ([]byte, error) {
	return buf, nil
}

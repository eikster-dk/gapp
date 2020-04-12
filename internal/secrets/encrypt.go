package secrets

import (
	"crypto/rand"
	"encoding/base64"
	"golang.org/x/crypto/nacl/box"
)

type Encrypt struct {

}

func NewEncrypt() Encrypt {
	return NewEncrypt()
}

func (e Encrypt) Encrypt(value string, pkey []byte) (string, error) {
	msg := []byte(value)
	var key [32]byte
	copy(key[:], pkey)

	var out []byte
	encrypted, err := box.SealAnonymous(out, msg, &key, rand.Reader)
	if err != nil {
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString(encrypted)

	return encoded, nil
}
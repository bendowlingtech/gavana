package gavana

import (
	"crypto/rand"
	"golang.org/x/crypto/argon2"
)

type Params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

func NewParams() *Params {
	p := &Params{
		memory:      64 * 1024,
		iterations:  3,
		parallelism: 2,
		saltLength:  16,
		keyLength:   32,
	}
	return p
}

func basicAuth() {

}

func generateFromPassword(password string, p *Params) {
	salt, err := generateRandomBytes(p.saltLength)
	if err != nil {
		return nil, err
	}
	hash := argon2.IDKey([]byte(password), )

}

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
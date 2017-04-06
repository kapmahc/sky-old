package security

import (
	"crypto/hmac"
	"crypto/sha512"
)

// NewHmac new hmac
func NewHmac(key []byte) *Hmac {
	return &Hmac{key: key}
}

// Hmac hmac
type Hmac struct {
	key []byte
}

// Sum sum hmac
func (p *Hmac) Sum(plain []byte) []byte {
	mac := hmac.New(sha512.New, p.key)
	return mac.Sum(plain)
}

// Chk chk hmac
func (p *Hmac) Chk(plain, code []byte) bool {
	return hmac.Equal(p.Sum(plain), code)
}

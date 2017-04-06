package security

import "crypto/aes"

// NewAes new aes-security
func NewAes(key []byte) (*Cipher, error) {
	cip, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return &Cipher{cip: cip}, nil
}

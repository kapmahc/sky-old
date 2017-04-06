package security

import "crypto/aes"

// NewAes new aes-security
func NewAes(key []byte) (*Factory, error) {
	cip, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return &Factory{cip: cip}, nil
}

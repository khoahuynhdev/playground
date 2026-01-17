package key

import (
	"crypto/rand"
	"crypto/rsa"
)

const defaultKeySize = 4096 // 4096 bits

// GenerateKeypair creates a new RSA key pair with the specified key size.
func GenerateKeypair(keySize int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	if keySize == 0 {
		keySize = defaultKeySize
	}
	privKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return nil, nil, err
	}
	return privKey, &privKey.PublicKey, nil
}

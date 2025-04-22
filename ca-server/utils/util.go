package utils

import (
	"crypto/rand"
	"math/big"
)

// NewSerialNum generates a random big integer suitable for use as a certificate serial number
func NewSerialNum() *big.Int {
	serialNumber, _ := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	return serialNumber
}

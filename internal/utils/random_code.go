package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func GenerateRandomCode() string {
	max := big.NewInt(1000000)
	num, err := rand.Int(rand.Reader, max)
	if err != nil {
		panic(err)
	}
	code := fmt.Sprintf("%06d", num.Int64())
	return code
}

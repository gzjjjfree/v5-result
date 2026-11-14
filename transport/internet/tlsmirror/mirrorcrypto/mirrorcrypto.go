package mirrorcrypto

import "github.com/gzjjjfree/v5-result/common/crypto"

//go:generate go run github.com/gzjjjfree/v5-result/common/errors/errorgen

func generateInitialAEADNonce() crypto.BytesGenerator {
	return crypto.GenerateIncreasingNonce([]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF})
}

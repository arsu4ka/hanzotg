package utils

import "regexp"

func IsValidTxHash(txHash string) bool {
	// Ethereum transaction hash format: 64 hexadecimal characters
	ethPattern := regexp.MustCompile(`^0x[0-9a-fA-F]{64}$`)
	// Tron transaction hash format: 64 hexadecimal characters
	tronPattern := regexp.MustCompile(`^[0-9a-fA-F]{64}$`)
	// Solana transaction hash format: 64 base58 characters
	solanaPattern := regexp.MustCompile(`^[0-9A-Za-z]{88}$`)

	return ethPattern.MatchString(txHash) || tronPattern.MatchString(txHash) || solanaPattern.MatchString(txHash)
}

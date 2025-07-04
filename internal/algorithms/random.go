package algorithms

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
)

func CryptoRandFloat64() (float64, error) {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		return 0, fmt.Errorf("failed to read random bytes: %w", err)
	}
	val := binary.BigEndian.Uint64(b)
	return float64(val>>11) / (1 << 53), nil
}

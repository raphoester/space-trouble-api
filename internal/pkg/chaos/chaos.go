package chaos

import (
	"crypto/sha256"
	"encoding/binary"
	"math/rand"
)

func stringToSeed(s string) uint64 {
	hash := sha256.Sum256([]byte(s))
	return binary.BigEndian.Uint64(hash[:8])
}

func Int(seed string, cap int) int {
	if cap == 0 {
		return 0
	}
	u := stringToSeed(seed)
	r := rand.New(rand.NewSource(int64(u)))
	return r.Int() % cap
}

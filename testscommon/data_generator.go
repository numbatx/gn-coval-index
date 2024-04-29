package testscommon

import (
	"math/big"
	"math/rand"
	"strconv"
)

// GenerateRandomFixedBytes generates a random byte slice of a predefined n size
func GenerateRandomFixedBytes(n int) []byte {
	ret := make([]byte, n)

	for i := 0; i < n; i++ {
		ret[i] = byte(rand.Int())
	}

	return ret
}

// GenerateRandomBytes generates a random byte slice
func GenerateRandomBytes() []byte {
	return []byte(strconv.Itoa(rand.Int()))
}

// GenerateRandomBigInt generates a random big.Int
func GenerateRandomBigInt() *big.Int {
	return big.NewInt(rand.Int63())
}

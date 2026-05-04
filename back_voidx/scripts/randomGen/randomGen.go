package randomGen

import (
	"crypto/rand"
	"math/big"
)

func GenerateString(size int, alphabet string) (string, error) {
	result := make([]byte, size)
	alphabetLen := big.NewInt(int64(len(alphabet)))
	for i := 0; i < size; i++ {
		idx, err := rand.Int(rand.Reader, alphabetLen)
		if err != nil {
			return "", err
		}
		result[i] = alphabet[idx.Int64()]
	}
	return string(result), nil
}

func GenerateUID() ([12]byte, error) {
	var uid [12]byte
	_, err := rand.Read(uid[:])
	return uid, err
}

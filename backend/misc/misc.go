package misc

import (
	"fmt"
	"os"
	"strconv"
)

const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
const base = uint64(len(alphabet))
const maxValue = uint64(14776336)
const encodedLen = 4

const multiplier uint64 = 6364136223846793005

var xorKey uint64

func Init() {
	secretStr := os.Getenv("SECRET_INT")
	secret, err := strconv.ParseInt(secretStr, 10, 64)
	if err != nil {
		panic(fmt.Errorf("invalid SECRET_INT: %w", err))
	}
	xorKey = uint64(secret)
}

func EncodeBase62(lastID int64) string {
	num := uint64(lastID)
	ob := ((num * multiplier) ^ xorKey) % maxValue
	buf := make([]byte, encodedLen)
	for i := encodedLen - 1; i >= 0; i-- {
		buf[i] = alphabet[ob%base]
		ob /= base
	}
	return string(buf)
}

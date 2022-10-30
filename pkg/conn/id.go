package conn

import (
	"encoding/base64"
	"math/big"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GenerateConnID() string {
	timeInt := big.NewInt(time.Now().UnixMilli())

	// 1951700038 epoch can represent as UJFORg
	time := base64.RawStdEncoding.EncodeToString(timeInt.Bytes())

	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, 4)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return time + string(b)
}

package generator

import (
	"crypto/rand"
	"fmt"
)

type RandomValueGenerator struct{}

func (generator *RandomValueGenerator) GetNextValue() string {
	delimitersMap := map[int]bool{
		4:  true,
		6:  true,
		8:  true,
		10: true,
	}

	randomBytesLength := 16
	randomBytes := make([]byte, randomBytesLength)
	_, err := rand.Read(randomBytes)

	if err != nil {
		panic(err)
	}

	uuid := ""

	for j := 0; j < randomBytesLength; j++ {
		if delimitersMap[j] {
			uuid += "-"
		}
		uuid += fmt.Sprintf("%02x", randomBytes[j])
	}

	return uuid
}

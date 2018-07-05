package service

import (
	"crypto/rand"
	"fmt"
	"net/http"
)

type Generator struct{}

type GenerateCommandArguments struct {
	Count    int    `json:"count"`
	Sequence string `json:"sequence"`
}

type GenerateResponse struct {
	Sequence string   `json:"sequence"`
	Ids      []string `json:"ids"`
}

func (generator *Generator) Generate(r *http.Request, args *GenerateCommandArguments, response *GenerateResponse) error {
	sequence := args.Sequence

	if sequence == "" {
		sequence = "default"
	}

	ids := make([]string, 0)

	for i := 0; i < args.Count; i++ {
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
			return err
		}

		uuid := ""

		for j := 0; j < randomBytesLength; j++ {
			if delimitersMap[j] {
				uuid += "-"
			}
			uuid += fmt.Sprintf("%02x", randomBytes[j])
		}

		ids = append(ids, uuid)
	}

	response.Sequence = sequence
	response.Ids = ids

	return nil
}

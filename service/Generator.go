package service

import (
	"github.com/strider2038/serial-uuid-generator/generator"
	"net/http"
)

type Generator struct {
	valueGenerator generator.ValueGenerator
}

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
	generator.valueGenerator.ReserveRange(sequence, uint64(args.Count))

	for i := 0; i < args.Count; i++ {
		ids = append(ids, generator.valueGenerator.GetNextValue(sequence))
	}

	response.Sequence = sequence
	response.Ids = ids

	return nil
}

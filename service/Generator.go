package service

import (
	"github.com/asaskevich/govalidator"
	"github.com/strider2038/serial-uuid-generator/generator"
	"net/http"
)

type Generator struct {
	valueGenerator generator.ValueGenerator
}

func NewGenerator(valueGenerator generator.ValueGenerator) *Generator {
	return &Generator{valueGenerator}
}

type GenerateCommandArguments struct {
	Count    int    `json:"count"    valid:"required,range(1|100000)"`
	Sequence string `json:"sequence" valid:"length(0|32),matches(^[A-Za-z\\-_]*$)"`
}

type GenerateResponse struct {
	Sequence string   `json:"sequence"`
	Ids      []string `json:"ids"`
}

func (generator *Generator) Generate(r *http.Request, args *GenerateCommandArguments, response *GenerateResponse) error {
	_, err := govalidator.ValidateStruct(args)

	if err != nil {
		return err
	}

	sequence := args.Sequence

	if sequence == "" {
		sequence = "default"
	}

	response.Sequence = sequence
	response.Ids = generator.generateIds(sequence, uint64(args.Count))

	return nil
}

func (generator *Generator) generateIds(sequence string, count uint64) []string {
	ids := make([]string, 0)
	generator.valueGenerator.ReserveRange(sequence, count)

	for i := uint64(0); i < count; i++ {
		ids = append(ids, generator.valueGenerator.GetNextValue(sequence))
	}

	return ids
}

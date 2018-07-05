package service

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGenerator_Generate_CountIsThreeAndSequenceIsEmptyString_SequenceIsDefaultAndThreeIdsGeneratedAndReturned(t *testing.T) {
	generator := Generator{}
	arguments := GenerateCommandArguments{3, ""}
	request := http.Request{}
	response := GenerateResponse{}

	generator.Generate(&request, &arguments, &response)

	assert.Equal(t, 3, len(response.Ids))
	assert.Equal(t, "default", response.Sequence)
}

func TestGenerator_Generate_CountIsOneAndSequenceIsCustom_SequenceIsCustomAndOneIdGeneratedAndReturned(t *testing.T) {
	generator := Generator{}
	arguments := GenerateCommandArguments{1, "custom"}
	request := http.Request{}
	response := GenerateResponse{}

	generator.Generate(&request, &arguments, &response)

	assert.Equal(t, 1, len(response.Ids))
	assert.Regexp(t, "[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}", response.Ids[0])
	assert.Equal(t, "custom", response.Sequence)
}

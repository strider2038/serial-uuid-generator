package service

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"
)

type valueGeneratorMock struct {
	mock.Mock
}

func (mock valueGeneratorMock) ReserveRange(sequence string, count uint64) {
	mock.Called(sequence, count)
}

func (mock valueGeneratorMock) GetNextValue(sequence string) string {
	args := mock.Called(sequence)

	return args.String(0)
}

func TestGenerator_Generate_CountIsThreeAndSequenceIsEmptyString_SequenceIsDefaultAndThreeIdsGeneratedAndReturned(t *testing.T) {
	const sequence = "default"
	const count = 3
	valueGenerator := valueGeneratorMock{}
	generator := Generator{&valueGenerator}
	arguments := GenerateCommandArguments{count, ""}
	request := http.Request{}
	response := GenerateResponse{}
	valueGenerator.On("GetNextValue", sequence).Return("value")
	valueGenerator.On("ReserveRange", sequence, uint64(count)).Return()

	generator.Generate(&request, &arguments, &response)

	valueGenerator.AssertExpectations(t)
	assert.Equal(t, count, len(response.Ids))
	assert.Equal(t, sequence, response.Sequence)
}

func TestGenerator_Generate_CountIsOneAndSequenceIsCustom_SequenceIsCustomAndOneIdGeneratedAndReturned(t *testing.T) {
	const sequence = "custom"
	const count = 1
	const value = "value"
	valueGenerator := valueGeneratorMock{}
	generator := Generator{&valueGenerator}
	arguments := GenerateCommandArguments{count, sequence}
	request := http.Request{}
	response := GenerateResponse{}
	valueGenerator.On("GetNextValue", sequence).Return(value)
	valueGenerator.On("ReserveRange", sequence, uint64(count)).Return()

	generator.Generate(&request, &arguments, &response)

	valueGenerator.AssertExpectations(t)
	assert.Equal(t, count, len(response.Ids))
	assert.Equal(t, value, response.Ids[0])
	assert.Equal(t, sequence, response.Sequence)
}

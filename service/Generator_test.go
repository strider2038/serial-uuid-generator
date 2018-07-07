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

func (mock valueGeneratorMock) GetNextValue() string {
	args := mock.Called()

	return args.String(0)
}

func TestGenerator_Generate_CountIsThreeAndSequenceIsEmptyString_SequenceIsDefaultAndThreeIdsGeneratedAndReturned(t *testing.T) {
	valueGenerator := valueGeneratorMock{}
	generator := Generator{&valueGenerator}
	arguments := GenerateCommandArguments{3, ""}
	request := http.Request{}
	response := GenerateResponse{}
	valueGenerator.On("GetNextValue").Return("value")

	generator.Generate(&request, &arguments, &response)

	valueGenerator.AssertExpectations(t)
	assert.Equal(t, 3, len(response.Ids))
	assert.Equal(t, "default", response.Sequence)
}

func TestGenerator_Generate_CountIsOneAndSequenceIsCustom_SequenceIsCustomAndOneIdGeneratedAndReturned(t *testing.T) {
	valueGenerator := valueGeneratorMock{}
	generator := Generator{&valueGenerator}
	arguments := GenerateCommandArguments{1, "custom"}
	request := http.Request{}
	response := GenerateResponse{}
	valueGenerator.On("GetNextValue").Return("value")

	generator.Generate(&request, &arguments, &response)

	valueGenerator.AssertExpectations(t)
	assert.Equal(t, 1, len(response.Ids))
	assert.Equal(t, "value", response.Ids[0])
	assert.Equal(t, "custom", response.Sequence)
}

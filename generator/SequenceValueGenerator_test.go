package generator

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type valueStorageMock struct {
	mock.Mock
}

func (mock *valueStorageMock) GetNextRangeForSequence(sequence string, count uint64) SequenceRange {
	args := mock.Called(sequence, count)

	return args.Get(0).(SequenceRange)
}

func TestSequenceValueGenerator_ReserveRange_NoRangeStepsInMapAndCount_RangeStepIsSetToUpperEdge(t *testing.T) {
	tests := []struct {
		name              string
		count             uint64
		reservedRangeStep uint64
	}{
		{
			"Count is 1 | Reserved range step not set",
			1,
			0,
		},
		{
			"Count is 100 | Reserved range step not set",
			100,
			0,
		},
		{
			"Count is 101 | Reserved range step set to 200",
			101,
			200,
		},
		{
			"Count is 199 | Reserved range step set to 200",
			199,
			200,
		},
		{
			"Count is 200 | Reserved range step set to 300",
			200,
			300,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			const sequence = "sequence"
			const rangeStep = 100
			valueStorage := valueStorageMock{}
			sequences := make(map[string]*SequenceRange)
			rangeStepsToReserve := make(map[string]uint64)
			valueGenerator := sequenceValueGenerator{
				&valueStorage,
				rangeStep,
				sequences,
				rangeStepsToReserve,
			}

			valueGenerator.ReserveRange(sequence, test.count)

			assert.Equal(t, uint64(test.reservedRangeStep), rangeStepsToReserve[sequence])
		})
	}
}

func TestSequenceValueGenerator_GetNextValue_NotInitializedSequence_RangeCreatedAndFirstValueReturned(t *testing.T) {
	const sequence = "sequence"
	const rangeStep = 100
	valueStorage := valueStorageMock{}
	sequences := make(map[string]*SequenceRange)
	rangeStepsToReserve := make(map[string]uint64)
	valueGenerator := sequenceValueGenerator{
		&valueStorage,
		rangeStep,
		sequences,
		rangeStepsToReserve,
	}
	valueStorage.On("GetNextRangeForSequence", sequence, uint64(rangeStep)).Return(SequenceRange{
		0,
		99,
	})

	value := valueGenerator.GetNextValue(sequence)

	valueStorage.AssertExpectations(t)
	assert.Regexp(t, "[0-9a-f]{8}-[0-9a-f]{4}-0[0-9a-f]{3}-0000-000000000001", value)
	assert.Equal(t, uint64(1), sequences[sequence].CurrentValue)
}

func TestSequenceValueGenerator_GetNextValue_InitialSequence_FirstValueReturned(t *testing.T) {
	const sequence = "sequence"
	const rangeStep = 100
	valueStorage := valueStorageMock{}
	sequences := make(map[string]*SequenceRange)
	sequences[sequence] = &SequenceRange{
		0,
		100,
	}
	rangeStepsToReserve := make(map[string]uint64)
	valueGenerator := sequenceValueGenerator{
		&valueStorage,
		rangeStep,
		sequences,
		rangeStepsToReserve,
	}

	value := valueGenerator.GetNextValue(sequence)

	assert.Regexp(t, "[0-9a-f]{8}-[0-9a-f]{4}-0[0-9a-f]{3}-0000-000000000001", value)
	assert.Equal(t, uint64(1), sequences[sequence].CurrentValue)
}

func TestSequenceValueGenerator_GetNextValue_SequenceRangeIsFull_NewRangeCreatedAndNextValueReturned(t *testing.T) {
	const sequence = "sequence"
	const rangeStep = 100
	valueStorage := valueStorageMock{}
	sequences := make(map[string]*SequenceRange)
	sequences[sequence] = &SequenceRange{
		99,
		99,
	}
	rangeStepsToReserve := make(map[string]uint64)
	valueGenerator := sequenceValueGenerator{
		&valueStorage,
		rangeStep,
		sequences,
		rangeStepsToReserve,
	}
	valueStorage.On("GetNextRangeForSequence", sequence, uint64(rangeStep)).Return(SequenceRange{
		200,
		300,
	})

	value := valueGenerator.GetNextValue(sequence)

	valueStorage.AssertExpectations(t)
	assert.Regexp(t, "[0-9a-f]{8}-[0-9a-f]{4}-0[0-9a-f]{3}-0000-0000000000c8", value)
	assert.Equal(t, uint64(200), sequences[sequence].CurrentValue)
}

func TestSequenceValueGenerator_GetNextValue_SequenceWithFullLowerPart_NextValueReturned(t *testing.T) {
	const sequence = "sequence"
	const rangeStep = 100
	valueStorage := valueStorageMock{}
	sequences := make(map[string]*SequenceRange)
	sequences[sequence] = &SequenceRange{
		0x0000FFFFFFFFFFFF,
		0xFFFFFFFFFFFFFFFF,
	}
	rangeStepsToReserve := make(map[string]uint64)
	valueGenerator := sequenceValueGenerator{
		&valueStorage,
		rangeStep,
		sequences,
		rangeStepsToReserve,
	}

	value := valueGenerator.GetNextValue(sequence)

	assert.Regexp(t, "[0-9a-f]{8}-[0-9a-f]{4}-0[0-9a-f]{3}-0001-000000000000", value)
	assert.Equal(t, uint64(0x1000000000000), sequences[sequence].CurrentValue)
}

func TestSequenceValueGenerator_GetNextValue_NotInitializedSequenceAndReservedRangeGreaterThanRangeStep_LongerRangeCreatedAndFirstValueReturned(t *testing.T) {
	const sequence = "sequence"
	const rangeStep = 100
	const reservedRangeStep = 200
	valueStorage := valueStorageMock{}
	sequences := make(map[string]*SequenceRange)
	rangeStepsToReserve := make(map[string]uint64)
	rangeStepsToReserve[sequence] = reservedRangeStep
	valueGenerator := sequenceValueGenerator{
		&valueStorage,
		rangeStep,
		sequences,
		rangeStepsToReserve,
	}
	valueStorage.On("GetNextRangeForSequence", sequence, uint64(reservedRangeStep)).Return(SequenceRange{
		0,
		199,
	})

	value := valueGenerator.GetNextValue(sequence)

	valueStorage.AssertExpectations(t)
	assert.Regexp(t, "[0-9a-f]{8}-[0-9a-f]{4}-0[0-9a-f]{3}-0000-000000000001", value)
	assert.Equal(t, uint64(1), sequences[sequence].CurrentValue)
	assert.Empty(t, rangeStepsToReserve[sequence])
}

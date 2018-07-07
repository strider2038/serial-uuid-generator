package generator

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

const (
	testSequence  = "testSequence"
	testRangeStep = 100
)

type valueStorageMock struct {
	mock.Mock
}

func (mock *valueStorageMock) GetNextRangeForSequence(sequence string, count uint64) SequenceRange {
	args := mock.Called(sequence, count)

	return args.Get(0).(SequenceRange)
}

type SequenceValueGeneratorSuite struct {
	suite.Suite
	valueStorage        valueStorageMock
	rangeStep           uint64
	sequences           map[string]*SequenceRange
	rangeStepsToReserve map[string]uint64
}

func (suite *SequenceValueGeneratorSuite) SetupTest() {
	suite.valueStorage = valueStorageMock{}
	suite.rangeStep = testRangeStep
	suite.sequences = make(map[string]*SequenceRange)
	suite.rangeStepsToReserve = make(map[string]uint64)
}

func (suite *SequenceValueGeneratorSuite) createSequenceValueGenerator() *sequenceValueGenerator {
	return &sequenceValueGenerator{
		&suite.valueStorage,
		suite.rangeStep,
		suite.sequences,
		suite.rangeStepsToReserve,
	}
}

func TestSequenceValueGeneratorSuite(t *testing.T) {
	suite.Run(t, new(SequenceValueGeneratorSuite))
}

func TestNewSequenceValueGenerator(t *testing.T) {
	valueStorage := valueStorageMock{}

	valueGenerator := NewSequenceValueGenerator(&valueStorage, testRangeStep)

	assert.NotNil(t, valueGenerator)
}

func (suite *SequenceValueGeneratorSuite) TestSequenceValueGenerator_ReserveRange_NoRangeStepsInMapAndCount_RangeStepIsSetToUpperEdge() {
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
		suite.T().Run(test.name, func(t *testing.T) {
			valueGenerator := suite.createSequenceValueGenerator()

			valueGenerator.ReserveRange(testSequence, test.count)

			assert.Equal(suite.T(), uint64(test.reservedRangeStep), suite.rangeStepsToReserve[testSequence])
		})
	}
}

func (suite *SequenceValueGeneratorSuite) TestSequenceValueGenerator_GetNextValue_NotInitializedSequence_RangeCreatedAndFirstValueReturned() {
	valueGenerator := suite.createSequenceValueGenerator()
	suite.valueStorage.On("GetNextRangeForSequence", testSequence, uint64(testRangeStep)).Return(SequenceRange{
		0,
		99,
	})

	value := valueGenerator.GetNextValue(testSequence)

	suite.valueStorage.AssertExpectations(suite.T())
	assert.Regexp(suite.T(), "[0-9a-f]{8}-[0-9a-f]{4}-0[0-9a-f]{3}-0000-000000000001", value)
	assert.Equal(suite.T(), uint64(1), suite.sequences[testSequence].CurrentValue)
}

func (suite *SequenceValueGeneratorSuite) TestSequenceValueGenerator_GetNextValue_RangeInitialized_NextValueReturned() {
	tests := []struct {
		name                        string
		currentValue                uint64
		expectedNextValue           uint64
		expectedReturnedValueRegExp string
	}{
		{
			"Initial sequence | Next value is 1",
			0,
			1,
			"[0-9a-f]{8}-[0-9a-f]{4}-0[0-9a-f]{3}-0000-000000000001",
		},
		{
			"Sequence with full lower part | Next value is valid",
			0x0000FFFFFFFFFFFF,
			0x0001000000000000,
			"[0-9a-f]{8}-[0-9a-f]{4}-0[0-9a-f]{3}-0001-000000000000",
		},
	}

	for _, test := range tests {
		suite.T().Run(test.name, func(t *testing.T) {
			suite.sequences[testSequence] = &SequenceRange{
				test.currentValue,
				0xFFFFFFFFFFFFFFFF,
			}
			valueGenerator := suite.createSequenceValueGenerator()

			value := valueGenerator.GetNextValue(testSequence)

			assert.Regexp(suite.T(), test.expectedReturnedValueRegExp, value)
			assert.Equal(suite.T(), test.expectedNextValue, suite.sequences[testSequence].CurrentValue)
		})
	}
}

func (suite *SequenceValueGeneratorSuite) TestSequenceValueGenerator_GetNextValue_SequenceRangeIsFull_NewRangeCreatedAndNextValueReturned() {
	suite.sequences[testSequence] = &SequenceRange{
		99,
		99,
	}
	valueGenerator := suite.createSequenceValueGenerator()
	suite.valueStorage.On("GetNextRangeForSequence", testSequence, uint64(testRangeStep)).Return(SequenceRange{
		200,
		300,
	})

	value := valueGenerator.GetNextValue(testSequence)

	suite.valueStorage.AssertExpectations(suite.T())
	assert.Regexp(suite.T(), "[0-9a-f]{8}-[0-9a-f]{4}-0[0-9a-f]{3}-0000-0000000000c8", value)
	assert.Equal(suite.T(), uint64(200), suite.sequences[testSequence].CurrentValue)
}

func (suite *SequenceValueGeneratorSuite) TestSequenceValueGenerator_GetNextValue_NotInitializedSequenceAndReservedRangeGreaterThanRangeStep_LongerRangeCreatedAndFirstValueReturned() {
	const reservedRangeStep = 200
	suite.rangeStepsToReserve[testSequence] = reservedRangeStep
	valueGenerator := suite.createSequenceValueGenerator()
	suite.valueStorage.On("GetNextRangeForSequence", testSequence, uint64(reservedRangeStep)).Return(SequenceRange{
		0,
		199,
	})

	value := valueGenerator.GetNextValue(testSequence)

	suite.valueStorage.AssertExpectations(suite.T())
	assert.Regexp(suite.T(), "[0-9a-f]{8}-[0-9a-f]{4}-0[0-9a-f]{3}-0000-000000000001", value)
	assert.Equal(suite.T(), uint64(1), suite.sequences[testSequence].CurrentValue)
	assert.Empty(suite.T(), suite.rangeStepsToReserve[testSequence])
}

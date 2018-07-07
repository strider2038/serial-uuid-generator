package generator

import (
	"crypto/rand"
	"fmt"
)

type sequenceValueGenerator struct {
	valueStorage        ValueStorage
	rangeStep           uint64
	sequences           map[string]*SequenceRange
	rangeStepsToReserve map[string]uint64
}

func NewSequenceValueGenerator(valueStorage ValueStorage, rangeStep uint64) ValueGenerator {
	sequences := make(map[string]*SequenceRange)
	rangeStepsToReserve := make(map[string]uint64)

	return &sequenceValueGenerator{
		valueStorage,
		rangeStep,
		sequences,
		rangeStepsToReserve,
	}
}

func (valueGenerator *sequenceValueGenerator) ReserveRange(sequence string, count uint64) {
	if count > valueGenerator.rangeStep {
		ratio := float64(count) / float64(valueGenerator.rangeStep)
		upperEdge := (uint64(ratio) + 1) * valueGenerator.rangeStep
		valueGenerator.rangeStepsToReserve[sequence] = upperEdge
	}
}

func (valueGenerator *sequenceValueGenerator) GetNextValue(sequence string) string {
	if valueGenerator.sequences[sequence] == nil {
		valueGenerator.updateRangeForSequence(sequence)
	}

	currentValue := &valueGenerator.sequences[sequence].CurrentValue
	*currentValue++

	if *currentValue > valueGenerator.sequences[sequence].MaxValue {
		currentValue = valueGenerator.updateRangeAndGetNextValueForSequence(sequence)
	}

	sequentialPart := valueGenerator.formatSequentialPart(currentValue)
	randomPart := valueGenerator.generateRandomPart(8)

	return randomPart + "-" + sequentialPart
}

func (valueGenerator *sequenceValueGenerator) updateRangeForSequence(sequence string) *SequenceRange {
	var rangeStep uint64

	if valueGenerator.rangeStepsToReserve[sequence] > 0 {
		rangeStep = valueGenerator.rangeStepsToReserve[sequence]
		delete(valueGenerator.rangeStepsToReserve, sequence)
	} else {
		rangeStep = valueGenerator.rangeStep
	}

	nextRange := valueGenerator.valueStorage.GetNextRangeForSequence(sequence, rangeStep)
	valueGenerator.sequences[sequence] = &nextRange

	return &nextRange
}

func (valueGenerator *sequenceValueGenerator) updateRangeAndGetNextValueForSequence(sequence string) *uint64 {
	delete(valueGenerator.sequences, sequence)
	nextRange := valueGenerator.updateRangeForSequence(sequence)

	return &nextRange.CurrentValue
}

func (valueGenerator *sequenceValueGenerator) formatSequentialPart(currentValue *uint64) string {
	lowerSequentialPartHex := fmt.Sprintf("%012x", *currentValue&0x0000FFFFFFFFFFFF)
	higherSequentialPartHex := fmt.Sprintf("%04x", *currentValue>>48&0x000000000000FFFF)
	sequentialPartHex := higherSequentialPartHex + "-" + lowerSequentialPartHex

	return sequentialPartHex
}

func (valueGenerator *sequenceValueGenerator) generateRandomPart(randomBytesCount int) string {
	randomBytes := make([]byte, randomBytesCount)
	_, err := rand.Read(randomBytes)

	if err != nil {
		panic(err)
	}

	randomHex := ""
	for i := 0; i < randomBytesCount; i++ {
		if i == 4 || i == 6 {
			randomHex += "-"
		}
		if i == 6 {
			// uuid version character special value = 0
			randomHex += fmt.Sprintf("0%01x", randomBytes[i]&0x0F)
		} else {
			randomHex += fmt.Sprintf("%02x", randomBytes[i])
		}
	}

	return randomHex
}

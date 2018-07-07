package generator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRandomValueGenerator_GetNextValue_NoParameters_UuidReturned(t *testing.T) {
	valueGenerator := RandomValueGenerator{}

	value := valueGenerator.GetNextValue("sequence")

	assert.Regexp(t, "[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}", value)
}

package config

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestLoadConfigFromEnvironment_AllParametersInEnv_ParametersLoaded(t *testing.T) {
	os.Setenv("UUIDGEN_DATABASE_URL", "postgres:5432")
	os.Setenv("UUIDGEN_TABLE_NAME", "public.table_name")
	os.Setenv("UUIDGEN_RANGE_STEP", "1000")
	os.Setenv("UUIDGEN_PORT", "5000")
	os.Setenv("UUIDGEN_LOG_LEVEL", "warn")
	os.Setenv("UUIDGEN_LOG_FORMAT", "json")

	config := LoadConfigFromEnvironment()

	assert.Equal(t, "postgres:5432", config.DatabaseUrl)
	assert.Equal(t, "public.table_name", config.TableName)
	assert.Equal(t, uint64(1000), config.RangeStep)
	assert.Equal(t, uint16(5000), config.Port)
	assert.Equal(t, logrus.WarnLevel, config.LogLevel)
	assert.Equal(t, "json", config.LogFormat)
}

func TestLoadConfigFromEnvironment_RequiredParametersInEnv_DefaultParametersLoaded(t *testing.T) {
	os.Setenv("UUIDGEN_DATABASE_URL", "postgres:5432")
	os.Unsetenv("UUIDGEN_TABLE_NAME")
	os.Unsetenv("UUIDGEN_RANGE_STEP")
	os.Unsetenv("UUIDGEN_PORT")
	os.Unsetenv("UUIDGEN_LOG_LEVEL")
	os.Unsetenv("UUIDGEN_LOG_FORMAT")

	config := LoadConfigFromEnvironment()

	assert.Equal(t, "postgres:5432", config.DatabaseUrl)
	assert.Equal(t, "public.uuid_sequence", config.TableName)
	assert.Equal(t, uint64(100), config.RangeStep)
	assert.Equal(t, uint16(3000), config.Port)
	assert.Equal(t, logrus.InfoLevel, config.LogLevel)
	assert.Equal(t, "tty", config.LogFormat)
}

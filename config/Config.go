package config

import (
	"os"
	"strconv"
)

type Config struct {
	DatabaseUrl string
	TableName   string
	RangeStep   uint64
	Port        uint16
}

func LoadConfigFromEnvironment() Config {
	databaseUrl := os.Getenv("UUIDGEN_DATABASE_URL")
	tableName := os.Getenv("UUIDGEN_TABLE_NAME")
	unparsedRangeStep := os.Getenv("UUIDGEN_RANGE_STEP")
	unparsedPort := os.Getenv("UUIDGEN_PORT")

	rangeStep, _ := strconv.ParseUint(unparsedRangeStep, 10, 64)
	port, _ := strconv.ParseUint(unparsedPort, 10, 32)

	if tableName == "" {
		tableName = "public.uuid_sequence"
	}

	if rangeStep == 0 {
		rangeStep = 100
	}

	if port == 0 {
		port = 3000
	}

	return Config{
		DatabaseUrl: databaseUrl,
		TableName:   tableName,
		RangeStep:   rangeStep,
		Port:        uint16(port),
	}
}

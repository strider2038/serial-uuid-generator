package config

import (
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
)

type Config struct {
	DatabaseUrl string
	TableName   string
	RangeStep   uint64
	Port        uint16
	LogLevel    logrus.Level
	LogFormat   string
}

func LoadConfigFromEnvironment() Config {
	return Config{
		DatabaseUrl: os.Getenv("UUIDGEN_DATABASE_URL"),
		TableName:   getTableName(),
		RangeStep:   getRangeStep(),
		Port:        getPort(),
		LogLevel:    getLogLevel(),
		LogFormat:   getLogFormat(),
	}
}

func getTableName() string {
	tableName := os.Getenv("UUIDGEN_TABLE_NAME")

	if tableName == "" {
		tableName = "public.uuid_sequence"
	}

	return tableName
}

func getRangeStep() uint64 {
	unparsedRangeStep := os.Getenv("UUIDGEN_RANGE_STEP")
	rangeStep, _ := strconv.ParseUint(unparsedRangeStep, 10, 64)

	if rangeStep == 0 {
		rangeStep = 100
	}

	return rangeStep
}

func getPort() uint16 {
	unparsedPort := os.Getenv("UUIDGEN_PORT")
	port, _ := strconv.ParseUint(unparsedPort, 10, 32)

	if port == 0 {
		port = 3000
	}

	return uint16(port)
}

func getLogLevel() logrus.Level {
	var logLevel logrus.Level
	var err error

	unparsedLogLevel := os.Getenv("UUIDGEN_LOG_LEVEL")

	if unparsedLogLevel == "" {
		logLevel = logrus.InfoLevel
	} else {
		logLevel, err = logrus.ParseLevel(unparsedLogLevel)
		panicOnError(err)
	}

	return logLevel
}

func getLogFormat() string {
	var format string

	if os.Getenv("UUIDGEN_LOG_FORMAT") == "json" {
		format = "json"
	} else {
		format = "tty"
	}

	return format
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

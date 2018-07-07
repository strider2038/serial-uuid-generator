package generator

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"testing"

	_ "github.com/lib/pq"
)

const (
	testConnectionString = "postgres://user:password@localhost/generator?sslmode=disable"
	testTableName        = "public.uuid_sequence"
)

func TestPostgresValueStorage_GetNextRangeForSequence_TableIsEmptyAndCountIs100_SequenceRangeReturned(t *testing.T) {
	truncateTable(testTableName)
	storage := NewPostgresValueStorage(testConnectionString, testTableName)

	sequenceRange := storage.GetNextRangeForSequence(testSequence, 100)

	assert.Equal(t, uint64(0), sequenceRange.CurrentValue)
	assert.Equal(t, uint64(99), sequenceRange.MaxValue)
}

func TestPostgresValueStorage_GetNextRangeForSequence_TableHasRangeAndCountIs150_NewSequenceRangeReturned(t *testing.T) {
	truncateTable(testTableName)
	storage := NewPostgresValueStorage(testConnectionString, testTableName)
	storage.GetNextRangeForSequence(testSequence, 100)

	sequenceRange := storage.GetNextRangeForSequence(testSequence, 150)

	assert.Equal(t, uint64(100), sequenceRange.CurrentValue)
	assert.Equal(t, uint64(249), sequenceRange.MaxValue)
}

func truncateTable(table string) {
	db, err := sql.Open("postgres", testConnectionString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec("TRUNCATE TABLE " + table)
	if err != nil {
		panic(err)
	}
}

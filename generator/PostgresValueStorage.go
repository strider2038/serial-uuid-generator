package generator

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type postgresValueStorage struct {
	connectionString string
	query            string
}

func NewPostgresValueStorage(connectionString string, tableName string) ValueStorage {
	query := `
WITH updated AS (
  INSERT INTO %s (sequence, start_value) VALUES ($1, $2)
  ON CONFLICT (sequence)
    DO
    UPDATE SET start_value = %s.start_value + EXCLUDED.start_value
  RETURNING sequence, start_value AS max_value
)
SELECT
  COALESCE(original.start_value, 0) AS start_value,
  (updated.max_value - 1) AS max_value
FROM %s original
RIGHT JOIN updated ON (original.sequence = updated.sequence)
`

	return &postgresValueStorage{
		connectionString,
		fmt.Sprintf(query, tableName, tableName, tableName),
	}
}

func (storage *postgresValueStorage) GetNextRangeForSequence(sequence string, count uint64) SequenceRange {
	db, err := sql.Open("postgres", storage.connectionString)
	storage.handleError(err)
	defer db.Close()

	sequenceRange := SequenceRange{}
	err = db.QueryRow(storage.query, sequence, count).Scan(&sequenceRange.CurrentValue, &sequenceRange.MaxValue)
	storage.handleError(err)

	return sequenceRange
}

func (storage *postgresValueStorage) handleError(err error) {
	if err != nil {
		panic(err)
	}
}

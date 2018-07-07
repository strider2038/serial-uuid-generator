package generator

type ValueStorage interface {
	GetNextRangeForSequence(sequence string, count uint64) SequenceRange
}

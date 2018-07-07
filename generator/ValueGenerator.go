package generator

type ValueGenerator interface {
	ReserveRange(sequence string, count uint64)
	GetNextValue(sequence string) string
}

package generator

type ValueGenerator interface {
	ReserveRange(sequence string, count int)
	GetNextValue(sequence string) string
}

package enum

type PrStatus int16

const (
	PrStatusOpen PrStatus = iota
	PrStatusMerged
)

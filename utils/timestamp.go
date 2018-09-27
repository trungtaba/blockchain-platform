package utils

const epochTimestamp = 946684800

// BlockTimestamp utility timestamp
type BlockTimestamp struct {
	Slot uint64
}

// NewBlockTimestamp ...
func NewBlockTimestamp(timestamp uint64, blockTime uint64) BlockTimestamp {
	slot := (timestamp - epochTimestamp) / blockTime
	return BlockTimestamp{Slot: slot}
}

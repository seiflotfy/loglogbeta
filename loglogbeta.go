package loglogbeta

import "errors"

// LogLogBeta is a sketch for cardinality estimation based on LogLog counting
type LogLogBeta struct {
	registers []uint8
	m         uint32
	precision uint8
	alpha     float64
}

// New returns a LogLogBeta sketch with 2^precision registers, where
// precision must be between 4 and 16
func New(precision uint8) (*LogLogBeta, error) {
	if precision > 16 || precision < 4 {
		return nil, errors.New("precision must be between 4 and 16")
	}
	m := uint32(1 << precision)
	return &LogLogBeta{
		m:         m,
		precision: precision,
		registers: make([]uint8, m),
		alpha:     alpha(float64(m)),
	}, nil
}

// NewDefault returns a LogLogBeta sketch with 2^14 registers
func NewDefault() *LogLogBeta {
	llb, _ := New(14)
	return llb
}

// Add inserts a value into the sketch
func (llb *LogLogBeta) Add(value []byte) {
	k, val := getPosVal(value, llb.precision)
	if llb.registers[k] < val {
		llb.registers[k] = val
	}
}

// Cardinality returns the number of unique elements added to the sketch
func (llb *LogLogBeta) Cardinality() uint64 {
	sum := regSum(llb.registers)
	ez := zeros(llb.registers)
	m := float64(llb.m)
	beta := beta(ez)
	return uint64(llb.alpha * m * (m - ez) / (beta + sum))
}

// Merge takes another LogLogBeta and combines it with llb one, making llb the union of both.
func (llb *LogLogBeta) Merge(other *LogLogBeta) error {
	if llb.precision != llb.precision {
		return errors.New("precisions must be equal")
	}

	for i, v := range llb.registers {
		if v < other.registers[i] {
			llb.registers[i] = other.registers[i]
		}
	}

	return nil
}

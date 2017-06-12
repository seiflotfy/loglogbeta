package loglogbeta

import (
	"errors"
	"time"
)

// RLogLogBeta is a sketch for retaining cardinality estimation based on LogLog counting
type RLogLogBeta struct {
	registers [][]int64
	m         uint32
	precision uint8
	alpha     float64
}

// NewRetaining returns a RLogLogBeta sketch with 2^precision registers, where
// precision must be between 4 and 16
func NewRetaining(precision uint8) (*RLogLogBeta, error) {
	if precision > 16 || precision < 4 {
		return nil, errors.New("precision must be between 4 and 16")
	}

	m := uint32(1 << precision)
	max := 64 - precision
	registers := make([][]int64, m, m)
	for i := range registers {
		registers[i] = make([]int64, max, max)
	}

	return &RLogLogBeta{
		m:         m,
		precision: precision,
		registers: registers,
		alpha:     alpha(float64(m)),
	}, nil
}

// NewDefaultRetaining returns a RLogLogBeta sketch with 2^14 registers
func NewDefaultRetaining() *RLogLogBeta {
	rllb, _ := NewRetaining(14)
	return rllb
}

// AddNow inserts a value into the sketch with the current timestamp
func (rllb *RLogLogBeta) AddNow(value []byte) {
	rllb.Add(value, time.Now())
}

// Add inserts a value into the sketch with given timestamp
func (rllb *RLogLogBeta) Add(value []byte, timestamp time.Time) {
	now := timestamp.Unix()
	k, val := getPosVal(value, rllb.precision)
	val--

	if rllb.registers[k][val] <= now {
		rllb.registers[k][val] = now
	}
}

// CardinalitySince returns the number of unique elements added to the sketch since a specific timestamp
func (rllb *RLogLogBeta) CardinalitySince(since time.Time) uint64 {
	m := float64(rllb.m)
	sum := regSumSince(rllb.registers, since)
	ez := zerosSince(rllb.registers, since.Unix())
	beta := beta(ez)
	return uint64(rllb.alpha * m * (m - ez) / (beta + sum))
}

// Cardinality returns the number of unique elements added to the sketch since the beginning of time
func (rllb *RLogLogBeta) Cardinality() uint64 {
	// FIXME: need to start with +1 seconds which is weird
	return rllb.CardinalitySince(time.Unix(1, 0))
}

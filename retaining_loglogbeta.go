package loglogbeta

import (
	"errors"
	"time"

	metro "github.com/dgryski/go-metro"
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
	max := 32 - precision
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
func NewDefaultRetaining(precision uint8) (*RLogLogBeta, error) {
	return NewRetaining(14)
}

// AddNow inserts a value into the sketch with the current timestamp
func (rllb *RLogLogBeta) AddNow(value []byte) {
	rllb.Add(value, time.Now())
}

// Add inserts a value into the sketch with given timestamp
func (rllb *RLogLogBeta) Add(value []byte, timestamp time.Time) {
	now := timestamp.Unix()
	x := metro.Hash64(value, 1337)
	max := 32 - rllb.precision
	val := rho(x<<(rllb.precision+32), max) - 1
	k := x >> uint(max+32)

	if rllb.registers[k][val] <= now {
		rllb.registers[k][val] = now
	}
}

// Cardinality returns the number of unique elements added to the sketch since a specific timestamp
func (rllb *RLogLogBeta) Cardinality(since time.Time) uint64 {
	m := float64(rllb.m)
	sum := regSumSince(rllb.registers, since)
	ez := zerosSince(rllb.registers, since.Unix())
	beta := beta(ez)
	return uint64(rllb.alpha * m * (m - ez) / (beta + sum))
}

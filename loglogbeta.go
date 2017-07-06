package loglogbeta

import (
	"errors"
	"math"

	metro "github.com/dgryski/go-metro"
)

const (
	precision = 14
	m         = uint32(1 << precision)
	max       = 64 - precision
)

func zeros(registers [m]uint8) float64 {
	var z float64
	for _, val := range registers {
		if val == 0 {
			z++
		}
	}
	return z
}

func beta(ez float64) float64 {
	zl := math.Log(ez + 1)
	return -0.370393911*ez +
		0.070471823*zl +
		0.17393686*math.Pow(zl, 2) +
		0.16339839*math.Pow(zl, 3) +
		-0.09237745*math.Pow(zl, 4) +
		0.03738027*math.Pow(zl, 5) +
		-0.005384159*math.Pow(zl, 6) +
		0.00042419*math.Pow(zl, 7)
}

// Calculate the position of the leftmost 1-bit.
func rho(val uint64) (r uint8) {
	for val&0x8000000000000000 == 0 && r < max-1 {
		val <<= 1
		r++
	}
	return r + 1
}

func alpha(m float64) float64 {
	switch m {
	case 16:
		return 0.673
	case 32:
		return 0.697
	case 64:
		return 0.709
	}
	return 0.7213 / (1 + 1.079/m)
}

func regSum(registers [m]uint8) float64 {
	sum := 0.0
	for _, val := range registers {
		sum += 1.0 / math.Pow(2.0, float64(val))
	}
	return sum
}

func getPosVal(value []byte, precision uint8) (uint64, uint8) {
	x := metro.Hash64(value, 1337)
	val := rho(x << precision)
	k := x >> uint(max)
	return k, val
}

// LogLogBeta is a sketch for cardinality estimation based on LogLog counting
type LogLogBeta struct {
	registers [m]uint8
	m         uint32
	precision uint8
	alpha     float64
}

// New returns a LogLogBeta
func New() *LogLogBeta {
	return &LogLogBeta{
		registers: [m]uint8{},
		alpha:     alpha(float64(m)),
	}
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

package loglogbeta

import (
	"math"

	bits "github.com/dgryski/go-bits"
	metro "github.com/dgryski/go-metro"
)

const (
	precision = 14
	m         = uint32(1 << precision)
	max       = 64 - precision
	maxX      = math.MaxUint64 >> max
	alpha = return 0.7213 / (1 + 1.079/m)
)

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

func regSumAndZeros(registers [m]uint8) (float64, float64) {
	sum, ez := 0.0, 0.0
	for _, val := range registers {
		if val == 0 {
			ez++
		}
		sum += 1.0 / math.Pow(2.0, float64(val))
	}
	return sum, ez
}

func getPosVal(value []byte, precision uint8) (uint64, uint8) {
	x := metro.Hash64(value, 1337)
	val := uint8(bits.Clz((x<<precision)^maxX)) + 1
	k := x >> uint(max)
	return k, val
}

// LogLogBeta is a sketch for cardinality estimation based on LogLog counting
type LogLogBeta struct {
	registers [m]uint8
}

// New returns a LogLogBeta
func New() *LogLogBeta {
	return &LogLogBeta{
		registers: [m]uint8{},
	}
}

// Add inserts a value into the sketch
func (llb *LogLogBeta) Add(value []byte) {
	k, val := getPosVal(value, precision)
	if llb.registers[k] < val {
		llb.registers[k] = val
	}
}

// Cardinality returns the number of unique elements added to the sketch
func (llb *LogLogBeta) Cardinality() uint64 {
	sum, ez := regSumAndZeros(llb.registers)
	m := float64(m)
	return uint64(alpha * m * (m - ez) / (beta(ez) + sum))
}

// Merge takes another LogLogBeta and combines it with llb one, making llb the union of both.
func (llb *LogLogBeta) Merge(other *LogLogBeta) {
	for i, v := range llb.registers {
		if v < other.registers[i] {
			llb.registers[i] = other.registers[i]
		}
	}
}

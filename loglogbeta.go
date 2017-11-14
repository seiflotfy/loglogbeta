package loglogbeta

import (
	"fmt"
	"math"
	bits "math/bits"

	metro "github.com/dgryski/go-metro"
)

const (
	precision = 14
	m         = uint32(1 << precision) // 16384
	max       = 64 - precision
	maxX      = math.MaxUint64 >> max
	alpha     = 0.7213 / (1 + 1.079/float64(m))
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

func regSumAndZeros(registers []uint8) (float64, float64) {
	sum, ez := 0.0, 0.0
	for _, val := range registers {
		if val == 0 {
			ez++
		}
		sum += 1.0 / math.Pow(2.0, float64(val))
	}
	return sum, ez
}

// LogLogBeta is a sketch for cardinality estimation based on LogLog counting
type LogLogBeta [m]uint8

// New returns a LogLogBeta
func New() *LogLogBeta {
	return new(LogLogBeta)
}

// AddHash takes in a "hashed" value (bring your own hashing)
func (llb *LogLogBeta) AddHash(x uint64) {
	k := x >> uint(max)
	val := uint8(bits.LeadingZeros64((x<<precision)^maxX)) + 1
	if llb[k] < val {
		llb[k] = val
	}
}

// Add inserts a value into the sketch
func (llb *LogLogBeta) Add(value []byte) {
	llb.AddHash(metro.Hash64(value, 1337))
}

// Cardinality returns the number of unique elements added to the sketch
func (llb *LogLogBeta) Cardinality() uint64 {
	sum, ez := regSumAndZeros(llb[:])
	m := float64(m)
	return uint64(alpha * m * (m - ez) / (beta(ez) + sum))
}

// Merge takes another LogLogBeta and combines it with llb one, making llb the union of both.
func (llb *LogLogBeta) Merge(other *LogLogBeta) {
	for i, v := range llb {
		if v < other[i] {
			llb[i] = other[i]
		}
	}
}

// Marshal returns a byte slice representation
func (llb *LogLogBeta) Marshal() []byte {
	return llb[:]
}

// Unmarshal returns a LogLogBeta for given bytes
func Unmarshal(bytes []byte) (*LogLogBeta, error) {
	if uint32(len(bytes)) != m {
		return nil, fmt.Errorf("Invalid byte slice length: expected %d, got %d", m, len(bytes))
	}
	llb := new(LogLogBeta)
	copy(llb[:], bytes)
	return llb, nil
}

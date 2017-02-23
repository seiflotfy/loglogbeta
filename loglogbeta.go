package loglogbeta

import (
	"errors"
	"math"

	metro "github.com/dgryski/go-metro"
)

const two32 = 1 << 32

func zeros(registers []uint8) float64 {
	var z float64
	for _, val := range registers {
		if val == 0 {
			z++
		}
	}
	return z
}

func beta(ez, zl float64) float64 {
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
func rho(val uint64, max uint8) uint8 {
	r := uint8(1)
	for val&0x80000000 == 0 && r <= max {
		val <<= 1
		r++
	}
	return r
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

// LogLogBeta ...
type LogLogBeta struct {
	registers []uint8
	m         uint32
	bits      uint8
	alpha     float64
}

// New ...
func New(precision uint32) (*LogLogBeta, error) {
	if precision > 16 || precision < 4 {
		return nil, errors.New("precision must be between 4 and 16")
	}
	m := uint32(1 << precision)
	return &LogLogBeta{
		m:         m,
		bits:      uint8(math.Ceil(math.Log2(float64(m)))),
		registers: make([]uint8, m),
		alpha:     alpha(float64(m)),
	}, nil
}

// Add ...
func (llb *LogLogBeta) Add(value []byte) {
	x := metro.Hash64(value, 1337)
	max := 64 - llb.bits
	val := rho(x<<llb.bits, max)
	k := x >> uint(max)

	if llb.registers[k] < val {
		llb.registers[k] = val
	}
}

func linearCounting(m float64, v float64) float64 {
	return m * math.Log(m/float64(v))
}

// Cardinality ...
func (llb *LogLogBeta) Cardinality() uint64 {
	sum := 0.0
	m := float64(llb.m)
	for _, val := range llb.registers {
		sum += 1.0 / math.Pow(2.0, float64(val))
	}

	ez := zeros(llb.registers)
	zl := math.Log(ez + 1)
	beta := beta(ez, zl)
	est := llb.alpha * m * (m - ez) / (beta + sum)

	return uint64(est)
}

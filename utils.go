package loglogbeta

import (
	"math"
	"time"

	metro "github.com/dgryski/go-metro"
)

func zeros(registers []uint8) float64 {
	var z float64
	for _, val := range registers {
		if val == 0 {
			z++
		}
	}
	return z
}

func zerosSince(registers [][]int64, since int64) float64 {
	var z float64
	for _, reg := range registers {
		isZero := true
		for _, val := range reg {
			if val >= since {
				isZero = false
				break
			}
		}
		if isZero {
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
func rho(val uint64, max uint8) (r uint8) {
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

func regSum(registers []uint8) float64 {
	sum := 0.0
	for _, val := range registers {
		sum += 1.0 / math.Pow(2.0, float64(val))
	}
	return sum
}

func regSumSince(registers [][]int64, since time.Time) float64 {
	sum := 0.0
	timestamp := since.Unix()
	for _, reg := range registers {
		maxI := 0
		for i := (len(reg) - 1); i >= 0; i-- {
			if reg[i] >= timestamp {
				maxI = i + 1
				break
			}
		}
		sum += 1.0 / math.Pow(2.0, float64(maxI))
	}
	return sum
}

func getPosVal(value []byte, precision uint8) (uint64, uint8) {
	x := metro.Hash64(value, 1337)
	max := 64 - precision
	val := rho(x<<(precision), max)
	k := x >> uint(max)
	return k, val
}

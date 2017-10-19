package loglogbeta

import (
	"math"
	"math/rand"
	"testing"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func estimateError(got, exp uint64) float64 {
	var delta uint64
	if got > exp {
		delta = got - exp
	} else {
		delta = exp - got
	}
	return float64(delta) / float64(exp)
}

/*
func TestZeros(t *testing.T) {
	registers := [m]uint8{}
	exp := 0.0
	for i := range registers {
		val := uint8(rand.Intn(32))
		if val == 0 {
			exp++
		}
		registers[i] = val
	}
	_, got := regSumAndZeros(registers)
	if got != exp {
		t.Errorf("expected %.2f, got %.2f", exp, got)
	}
}
*/

func RandStringBytesMaskImprSrc(n uint32) string {
	b := make([]byte, n)
	for i := uint32(0); i < n; i++ {
		b[i] = letterBytes[rand.Int()%len(letterBytes)]
	}
	return string(b)
}

func TestCardinality(t *testing.T) {
	llb := New()
	step := 10000
	unique := map[string]bool{}

	for i := 1; len(unique) <= 1000000; i++ {
		str := RandStringBytesMaskImprSrc(rand.Uint32() % 32)
		llb.Add([]byte(str))
		unique[str] = true

		if len(unique)%step == 0 {
			exact := uint64(len(unique))
			res := uint64(llb.Cardinality())
			step *= 10

			ratio := 100 * estimateError(res, exact)
			if ratio > 2 {
				t.Errorf("Exact %d, got %d which is %.2f%% error", exact, res, ratio)
			}

		}
	}
}

func TestMerge(t *testing.T) {
	llb1 := New()
	llb2 := New()

	unique := map[string]bool{}

	for i := 1; i <= 3500000; i++ {
		str := RandStringBytesMaskImprSrc(rand.Uint32() % 32)
		llb1.Add([]byte(str))
		unique[str] = true

		str = RandStringBytesMaskImprSrc(rand.Uint32() % 32)
		llb2.Add([]byte(str))
		unique[str] = true
	}

	llb1.Merge(llb2)
	exact := len(unique)
	res := int(llb1.Cardinality())

	ratio := 100 * math.Abs(float64(res-exact)) / float64(exact)
	expectedError := 1.04 / math.Sqrt(float64(m))

	if float64(res) < float64(exact)-(float64(exact)*expectedError) || float64(res) > float64(exact)+(float64(exact)*expectedError) {
		t.Errorf("Exact %d, got %d which is %.2f%% error", exact, res, ratio)
	}

	llb1.Merge(llb2)
	exact = res
	res = int(llb1.Cardinality())

	if float64(res) < float64(exact)-(float64(exact)*expectedError) || float64(res) > float64(exact)+(float64(exact)*expectedError) {
		t.Errorf("Exact %d, got %d which is %.2f%% error", exact, res, ratio)
	}
}

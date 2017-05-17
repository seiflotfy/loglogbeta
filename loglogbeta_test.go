package loglogbeta

import (
	"math"
	"math/rand"
	"testing"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var src = rand.NewSource(time.Now().UnixNano())

func RandStringBytesMaskImprSrc(n uint32) string {
	b := make([]byte, n)
	for i := uint32(0); i < n; i++ {
		b[i] = letterBytes[rand.Int()%len(letterBytes)]
	}
	return string(b)
}

func TestCardinality(t *testing.T) {
	llb, err := New(14)
	if err != nil {
		t.Error("expected no error, got", err)
	}

	step := 10000
	unique := map[string]bool{}

	for i := 1; len(unique) <= 1000000; i++ {
		str := RandStringBytesMaskImprSrc(rand.Uint32() % 32)
		llb.Add([]byte(str))
		unique[str] = true

		if len(unique)%step == 0 {
			exact := len(unique)
			step *= 10
			res := int(llb.Cardinality())
			ratio := 100 * math.Abs(float64(res-exact)) / float64(exact)

			expectedError := 1.04 / math.Sqrt(float64(llb.m))

			if float64(res) < float64(exact)-(float64(exact)*expectedError) || float64(res) > float64(exact)+(float64(exact)*expectedError) {
				t.Errorf("Exact %d, got %d which is %.2f%% error", exact, res, ratio)
			}
		}
	}
}

func TestMerge(t *testing.T) {
	llb1, err := New(14)
	if err != nil {
		t.Error("expected no error, got", err)
	}
	llb2, err := New(14)
	if err != nil {
		t.Error("expected no error, got", err)
	}

	unique := map[string]bool{}

	for i := 1; i <= 100000; i++ {
		str := RandStringBytesMaskImprSrc(rand.Uint32() % 32)
		llb1.Add([]byte(str))
		unique[str] = true

		str = RandStringBytesMaskImprSrc(rand.Uint32() % 32)
		llb2.Add([]byte(str))
		unique[str] = true
	}

	err = llb1.Merge(llb2)
	if err != nil {
		t.Error("expected no error, got", err)
	}

	exact := len(unique)
	res := int(llb1.Cardinality())
	ratio := 100 * math.Abs(float64(res-exact)) / float64(exact)
	expectedError := 1.04 / math.Sqrt(float64(llb1.m))

	if float64(res) < float64(exact)-(float64(exact)*expectedError) || float64(res) > float64(exact)+(float64(exact)*expectedError) {
		t.Errorf("Exact %d, got %d which is %.2f%% error", exact, res, ratio)
	}

	err = llb1.Merge(llb2)
	if err != nil {
		t.Error("expected no error, got", err)
	}

	exact = res
	res = int(llb1.Cardinality())

	if float64(res) < float64(exact)-(float64(exact)*expectedError) || float64(res) > float64(exact)+(float64(exact)*expectedError) {
		t.Errorf("Exact %d, got %d which is %.2f%% error", exact, res, ratio)
	}

}

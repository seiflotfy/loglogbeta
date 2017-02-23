package loglogbeta

import (
	"math/rand"
	"testing"
	"time"

	"math"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var src = rand.NewSource(time.Now().UnixNano())

func RandStringBytesMaskImprSrc(n uint32) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i := uint32(0); i < n; i++ {
		b[i] = letterBytes[rand.Int()%len(letterBytes)]
	}

	return string(b)
}
func TestP(t *testing.T) {
	step := 10000
	llb, err := New(14)
	if err != nil {
		t.Error("expected no error, got", err)
	}
	for i := 1; i <= 10000000; i++ {
		str := RandStringBytesMaskImprSrc(rand.Uint32() % 32)
		llb.Add([]byte(str))
		if i%step == 0 {
			step *= 10
			res := int(llb.Cardinality())
			ratio := 100 * math.Abs(float64(res-i)) / float64(i)
			if ratio > 5 {
				t.Errorf("Exact %d, got %d which is %.2f%% error", i, res, ratio)
			}
		}
	}
}

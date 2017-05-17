package loglogbeta

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"
)

func TestRetainingCardinalitySince(t *testing.T) {
	rllb, err := NewRetaining(14)
	if err != nil {
		t.Error("expected no error, got", err)
	}

	step := 10000
	unique := map[string]bool{}
	tmpUnique := map[string]bool{}
	since := []time.Time{time.Now()}

	for i := 1; len(unique) <= 1000000; i++ {
		str := RandStringBytesMaskImprSrc(rand.Uint32() % 32)
		timestamp := time.Now()
		rllb.Add([]byte(str), timestamp)
		unique[str] = true
		tmpUnique[str] = true
		if len(tmpUnique)%step == 0 {
			step *= 10
			now := time.Now()
			since = append(since, now)
			exact := len(tmpUnique)
			res := int(rllb.CardinalitySince(now))
			ratio := 100 * math.Abs(float64(res-exact)) / float64(exact)
			expectedError := 1.04 / math.Sqrt(float64(rllb.m))

			if float64(res) < float64(exact)-(float64(exact)*expectedError) || float64(res) > float64(exact)+(float64(exact)*expectedError) {
				t.Errorf("Exact %d, got %d which is %.2f%% error", exact, res, ratio)
			}

			exact = len(unique)
			res = int(rllb.CardinalitySince(since[0]))
			ratio = 100 * math.Abs(float64(res-exact)) / float64(exact)
			expectedError = 1.04 / math.Sqrt(float64(rllb.m))

			if float64(res) < float64(exact)-(float64(exact)*expectedError) || float64(res) > float64(exact)+(float64(exact)*expectedError) {
				t.Errorf("Exact %d, got %d which is %.2f%% error", exact, res, ratio)
			}

			tmpUnique = map[string]bool{}
			time.Sleep(time.Second * 1)
		}
	}
}

func TestRetainingLogLogBetaManual(t *testing.T) {
	rllb, err := NewRetaining(14)
	if err != nil {
		t.Error("expected no error, got", err)
	}

	for i := 1; i <= 1000; i++ {
		rllb.AddNow([]byte(fmt.Sprintf("flow-%d", i)))
	}

	timestamp1 := time.Now()
	exact := 1000
	res := int(rllb.CardinalitySince(timestamp1))
	ratio := 100 * math.Abs(float64(res-exact)) / float64(exact)
	expectedError := 1.04 / math.Sqrt(float64(rllb.m))

	if float64(res) < float64(exact)-(float64(exact)*expectedError) || float64(res) > float64(exact)+(float64(exact)*expectedError) {
		t.Errorf("Exact %d, got %d which is %.2f%% error", exact, res, ratio)
	}

	time.Sleep(time.Second)

	for i := 1001; i <= 2000; i++ {
		rllb.AddNow([]byte(fmt.Sprintf("flow-%d", i)))
	}

	timestamp2 := time.Now()
	exact = 1000
	res = int(rllb.CardinalitySince(timestamp2))
	ratio = 100 * math.Abs(float64(res-exact)) / float64(exact)
	expectedError = 1.04 / math.Sqrt(float64(rllb.m))

	if float64(res) < float64(exact)-(float64(exact)*expectedError) || float64(res) > float64(exact)+(float64(exact)*expectedError) {
		t.Errorf("Exact %d, got %d which is %.2f%% error", exact, res, ratio)
	}

	exact = 2000
	res = int(rllb.CardinalitySince(timestamp1))
	ratio = 100 * math.Abs(float64(res-exact)) / float64(exact)
	expectedError = 1.04 / math.Sqrt(float64(rllb.m))

	if float64(res) < float64(exact)-(float64(exact)*expectedError) || float64(res) > float64(exact)+(float64(exact)*expectedError) {
		t.Errorf("Exact %d, got %d which is %.2f%% error", exact, res, ratio)
	}

	time.Sleep(time.Second)

	for i := 2001; i <= 3000; i++ {
		rllb.AddNow([]byte(fmt.Sprintf("flow-%d", i)))
	}

	timestamp3 := time.Now()
	exact = 1000
	res = int(rllb.CardinalitySince(timestamp3))
	ratio = 100 * math.Abs(float64(res-exact)) / float64(exact)
	expectedError = 1.04 / math.Sqrt(float64(rllb.m))

	if float64(res) < float64(exact)-(float64(exact)*expectedError) || float64(res) > float64(exact)+(float64(exact)*expectedError) {
		t.Errorf("Exact %d, got %d which is %.2f%% error", exact, res, ratio)
	}

	exact = 2000
	res = int(rllb.CardinalitySince(timestamp2))
	ratio = 100 * math.Abs(float64(res-exact)) / float64(exact)
	expectedError = 1.04 / math.Sqrt(float64(rllb.m))

	if float64(res) < float64(exact)-(float64(exact)*expectedError) || float64(res) > float64(exact)+(float64(exact)*expectedError) {
		t.Errorf("Exact %d, got %d which is %.2f%% error", exact, res, ratio)
	}

	exact = 3000
	res = int(rllb.CardinalitySince(timestamp1))
	ratio = 100 * math.Abs(float64(res-exact)) / float64(exact)
	expectedError = 1.04 / math.Sqrt(float64(rllb.m))

	if float64(res) < float64(exact)-(float64(exact)*expectedError) || float64(res) > float64(exact)+(float64(exact)*expectedError) {
		t.Errorf("Exact %d, got %d which is %.2f%% error", exact, res, ratio)
	}

	time.Sleep(time.Second)
	timestamp4 := time.Now()

	exact = 0
	res = int(rllb.CardinalitySince(timestamp4))
	ratio = 100 * math.Abs(float64(res-exact)) / float64(exact)
	expectedError = 1.04 / math.Sqrt(float64(rllb.m))

	if float64(res) < float64(exact)-(float64(exact)*expectedError) || float64(res) > float64(exact)+(float64(exact)*expectedError) {
		t.Errorf("Exact %d, got %d which is %.2f%% error", exact, res, ratio)
	}

	exact = 3000
	res = int(rllb.Cardinality())
	ratio = 100 * math.Abs(float64(res-exact)) / float64(exact)
	expectedError = 1.04 / math.Sqrt(float64(rllb.m))

	if float64(res) < float64(exact)-(float64(exact)*expectedError) || float64(res) > float64(exact)+(float64(exact)*expectedError) {
		t.Errorf("Exact %d, got %d which is %.2f%% error", exact, res, ratio)
	}
}

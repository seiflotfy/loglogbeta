package loglogbeta

import (
	"fmt"
	"testing"

	metro "github.com/dgryski/go-metro"
)

const testm = 4

func TestM(t *testing.T) {

	M := [2][testm]uint8{
		[testm]uint8{1, 2, 3, 4},
		[testm]uint8{4, 3, 2, 1},
	}

	var ns [testm]uint8
	for i := range ns {
		ns[i] = max
	}

	for _, llb := range M {
		for j, val := range llb {
			if val < ns[j] {
				ns[j] = val
			}
		}
	}

	fmt.Println(ns)
}

func TestIntersection(t *testing.T) {
	llb1 := New()
	llb2 := New()
	llb3 := New()
	llb4 := New()
	llb5 := New()
	llb6 := New()
	unique := map[string]bool{}
	total := 0

	for i := 1; len(unique) <= 1000000; i++ {
		str := fmt.Sprintf("stream-%d", i)
		unique[str] = true

		x := metro.Hash64([]byte(str), 1337)
		k, val := getPosVal(x)

		if i%1 == 0 {
			llb1.add(k, val)
		}
		if i%2 == 0 {
			llb2.add(k, val)
		}
		if i%3 == 0 {
			llb3.add(k, val)
		}
		if i%4 == 0 {
			llb4.add(k, val)
		}
		if i%5 == 0 {
			llb5.add(k, val)
		}
		if i%6 == 0 {
			llb6.add(k, val)
		}
		if i%5 == 0 && i%6 == 0 {
			total++
		}
	}
	joint([]*LogLogBeta{llb5, llb6})
	fmt.Println("joint card (actual):   ", total)
}

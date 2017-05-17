package loglogbeta

import "testing"

func TestZeros(t *testing.T) {
	registers := []uint8{
		0, 1, 0, 5, 6, 7, 0, 9, 10, 11,
	}
	got := zeros(registers)
	exp := 3.0
	if got != exp {
		t.Errorf("expected %.2f, got %.2f", exp, got)
	}
}

func TestZerosSince(t *testing.T) {
	registers := [][]int64{
		[]int64{0, 1, 2, 4, 5, 6},
		[]int64{0, 1, 2, 4, 5, 6},
		[]int64{0, 1, 2, 4, 5, 7},
		[]int64{0, 1, 2, 4, 10, 12},
		[]int64{0, 1, 2, 4, 5, 20},
	}
	got := zerosSince(registers, 20)
	exp := 4.0
	if got != exp {
		t.Errorf("expected %.2f, got %.2f", exp, got)
	}
	got = zerosSince(registers, 12)
	exp = 3.0
	if got != exp {
		t.Errorf("expected %.2f, got %.2f", exp, got)
	}
	got = zerosSince(registers, 11)
	exp = 3.0
	if got != exp {
		t.Errorf("expected %.2f, got %.2f", exp, got)
	}
	got = zerosSince(registers, 4)
	exp = 0.0
	if got != exp {
		t.Errorf("expected %.2f, got %.2f", exp, got)
	}
}

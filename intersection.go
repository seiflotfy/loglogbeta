package loglogbeta

import (
	"fmt"
	"math"
)

func ln(n float64) float64 {
	var lo, hi, m float64
	if n <= 0 {
		return -1
	}

	if n == 1 {
		return 0
	}

	EPS := 0.00001
	lo = 0
	hi = n

	for math.Abs(lo-hi) >= EPS {
		m = float64((lo + hi) / 2.0)
		if math.Exp(m)-n < 0 {
			lo = m
		} else {
			hi = m
		}
	}
	return float64((lo + hi) / 2.0)
}

func generation(ns []uint64, nstar, k float64) float64 {
	s := float64(m)
	mul := 1.0
	for _, nj := range ns {
		pow := (float64(nj) - nstar) / (s * math.Pow(2, k))
		mul *= (1.0 - math.Exp(-pow))
	}
	return math.Exp(-nstar/(s*math.Pow(2, k))) * (1 - mul)
}

func dGeneration(ns []uint64, n, k float64) float64 {
	s := float64(m)
	s2k := (s * math.Pow(2, k))
	left := (1 / s2k) * math.Exp(-n/s2k)
	middle := 1.0
	right := 1.0
	for _, nj := range ns {
		pow := (float64(nj) - n) / s2k
		//fmt.Println(">", pow)
		middle += 1 / (math.Exp(pow) - 1)
		right *= (1.0 - math.Exp(-pow))
	}
	//fmt.Println(left, middle, right)
	return left * ((middle * right) - 1)
}

func probability(ns []uint64, n uint64, k uint8) float64 {
	switch {
	case k == 0:
		return generation(ns, float64(n), float64(k))
	case 0 < k && k < max:
		return generation(ns, float64(n), float64(k)) - generation(ns, float64(n), float64(k-1))
	case k == max:
		return 1 - generation(ns, float64(n), float64(k-1))
	}
	return 0.0
}

func dProbability(ns []uint64, n uint64, k uint8) float64 {
	switch {
	case k == 0:
		return dGeneration(ns, float64(n), float64(k))
	case 0 < k && k < max:
		return dGeneration(ns, float64(n), float64(k)) - dGeneration(ns, float64(n), float64(k-1))
	case k == max:
		return -dGeneration(ns, float64(n), float64(k-1))
	}
	return 0.0
}

func countDist(ns [m]uint8) []uint8 {
	res := make([]uint8, max+1)
	for _, val := range ns {
		res[val]++
	}
	return res
}

func intersectedM(M []*LogLogBeta) *LogLogBeta {
	var ns [m]uint8
	for i := range ns {
		ns[i] = max
	}
	for _, llb := range M {
		for j, val := range llb.registers {
			if val < ns[j] {
				ns[j] = val
			}
		}
	}
	return &LogLogBeta{
		registers: ns,
	}
}

func joint(llbs []*LogLogBeta) float64 {
	var (
		res     = 0.0
		M       = intersectedM(llbs)
		Nks     = countDist(M.registers)
		ns      = make([]uint64, len(llbs))
		maxCard = M.Cardinality()
		maxSum  = 0.0
		maxN    = uint64(0)
	)
	for n := uint64(1); n <= maxCard; n++ {
		sum := ln(alpha)
		for k, Nk := range Nks {
			p := probability(ns, n, uint8(k))
			sum += (float64(Nk) * ln(p))
		}
		if maxSum < sum {
			maxSum = sum
			maxN = n
		}
		fmt.Println(">>>", n, sum, maxCard)
	}
	fmt.Println(maxN, maxSum)
	return res
}

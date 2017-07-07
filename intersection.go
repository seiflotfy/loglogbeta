package loglogbeta

import "math"

func generation(ns [m]uint8, n, k float64) float64 {
	s := float64(len(ns))
	mul := 1.0
	for _, nj := range ns {
		pow := (float64(nj) - n) / (s * math.Pow(2, k))
		mul *= (1.0 - math.Exp(-pow))
	}
	return math.Exp(-n/(s*math.Pow(2, k))) * (1 - mul)
}

func dGeneration(ns [m]uint8, n, k float64) float64 {
	s := float64(len(ns))
	left := (1 / (s * math.Pow(2, k))) * math.Exp(-n/(s*math.Pow(2, k)))
	middle := 1.0
	right := 1.0
	for _, nj := range ns {
		pow := (float64(nj) - n) / (s * math.Pow(2, k))
		middle += math.Pow((math.Exp(pow) - 1), -1)
		right *= (1.0 - math.Exp(-pow))
	}
	return left * ((middle * right) - 1)
}

func probability(ns [m]uint8, n uint64, k, h uint8) float64 {
	switch {
	case k == 0:
		return generation(ns, float64(n), float64(k))
	case 0 < k && k < h:
		return generation(ns, float64(n), float64(k)) - generation(ns, float64(n), float64(k-1))
	case k == h:
		return 1 - generation(ns, float64(n), float64(k-1))
	}
	return 0.0
}

func dProbability(ns [m]uint8, n uint64, k, h uint8) float64 {
	switch {
	case k == 0:
		return dGeneration(ns, float64(n), float64(k))
	case 0 < k && k < h:
		return dGeneration(ns, float64(n), float64(k)) - dGeneration(ns, float64(n), float64(k-1))
	case k == h:
		return -dGeneration(ns, float64(n), float64(k-1))
	}
	return 0.0
}

func countDist(ns [m]uint8, h uint8) []uint8 {
	res := make([]uint8, h+1)
	for _, val := range ns {
		res[val]++
	}
	return res
}

func minM(M [][m]uint8) []uint8 {
	ns := M[0]
	for i, m := range M {
		for j, val := range m {
			if val < ns[j] {
				ns[j] = val
			}
		}
	}
	return ns
}

func joint(h uint8, M [][m]uint8) float64 {
	ns := minM(M)
	nk := countDist(ns, h)
	sum := 0.0
	tllb := LogLogBeta{
		registers: ns,
	}
	for n := 0; i < tllb.Cardinality(); n++ {
		for k, val := range nk {
			sum += (float64(val) * dProbability(ns, n, uint8(k), h) / probability(ns, n, uint8(k), h))
		}
	}
	return sum
}

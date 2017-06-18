# LogLog-Beta [![GoDoc](https://godoc.org/github.com/seiflotfy/loglogbeta?status.svg)](https://godoc.org/github.com/seiflotfy/loglogbeta)

[LogLog-Beta and More: A New Algorithm for Cardinality Estimation Based on LogLog Counting](https://arxiv.org/pdf/1612.02284.pdf) -
by Jason Qin, Denys Kim, Yumei Tung

**TL;DR:**
Better than HyperLogLog in approximating the number unique elements in a set

## LogLog-Beta

LogLog-Beta is a new algorithm for estimating cardinalities based on LogLog counting. The new algorithm uses only one formula and needs no additional bias corrections for the entire range of cardinalities, therefore, it is more efficient and simpler to implement. The simulations show that the accuracy provided by the new algorithm is as good as or better than the accuracy provided by either of HyperLogLog or HyperLogLog++.

#### Using LogLog-Beta

```go
// Create LogLogBeta
llb := loglogbeta.NewDefault()

// Add value to LogLogBeta
llb.Add([]byte("hello"))

// Returns cardinality
llb.Cardinality()
```

## Initial Results

From [demo](llbdemo/main.go)

```bash
file:  data/words-1
exact: 150
estimate: 157
ratio: 4.458599%

file:  data/words-2
exact: 1308
estimate: 1373
ratio: 4.734159%

file:  data/words-3
exact: 76205
estimate: 78440
ratio: 2.849312%

file:  data/words-4
exact: 235886
estimate: 236587
ratio: 0.296297%

file:  data/words-5
exact: 349900
estimate: 349388
ratio: -0.146542%

file:  data/words-6
exact: 479829
estimate: 481224
ratio: 0.289886%

total
exact: 660131
estimate: 660548
ratio: 0.063129%
 ```

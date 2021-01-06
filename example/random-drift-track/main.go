package main

import (
	"github.com/linger1216/go-utils/convert"
	"github.com/linger1216/go-utils/inout"
	"math/rand"
	"time"
)

/*
拟合方程式：Y = a + b·X + c·X2
参数：
a = -0.428571428571438
b = 1.67619047619048
c = -0.0476190476190482
*/

var (
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func main() {
	outputFilename := "res/drift_random.csv"
	outputFileCsv := inout.NewCsvFile(inout.NewCsvFileConfig("trunc", outputFilename, 4096))

	size := 100
	driftCount := 10

	driftIndexs := r.Perm(size)[:driftCount]

	_in := func(v int, arr []int) bool {
		for i := range arr {
			if arr[i] == v {
				return true
			}
		}
		return false
	}

	for i := 0; i < size; i++ {
		x := i
		y := x + r.Intn(10)
		if _in(i, driftIndexs) {
			sign := r.Int()%2 == 0
			offset := r.Intn(size / 2)
			if sign {
				y += offset
			} else {
				y -= offset
			}
		}
		_ = outputFileCsv.Exec(convert.Int64ToString(int64(x)), convert.Int64ToString(int64(y)))
	}
	_ = outputFileCsv.Close()
}

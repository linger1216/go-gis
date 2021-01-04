package main

import (
	"encoding/csv"
	"github.com/linger1216/go-utils/convert"
	"github.com/linger1216/go-utils/inout"
	"io"
	"math/rand"
	"os"
	"time"
)

var (
	r   = rand.New(rand.NewSource(time.Now().UnixNano()))
	Max = 0.00005
	Min = 0.00004
)

func randomBuffer(x float64) float64 {
	sign := r.Int()%2 == 0
	if sign {
		x += r.Float64()*(Max-Min) + Min
	} else {
		x -= r.Float64()*(Max-Min) + Min
	}
	return x
}

func main() {
	count := 0
	jump := 1
	rawFilename := "res/test2_line_RaPts1000.csv"
	rawLatIndex := 3
	rawLngIndex := 2
	outputFilename := "res/test2_line_RaPts1000_random.csv"
	outputFileCsv := inout.NewCsvFile(inout.NewCsvFileConfig("trunc", outputFilename, 4096))

	// load rawView
	rawFile, err := os.OpenFile(rawFilename, os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	rawReader := csv.NewReader(rawFile)
	rawReader.LazyQuotes = true
	for {
		line, err := rawReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		count++
		if count <= jump {
			continue
		}
		if len(line) >= 0 {
			Latitude := convert.StringToFloat(line[rawLatIndex])
			Longitude := convert.StringToFloat(line[rawLngIndex])

			//last := Latitude
			Latitude = randomBuffer(Latitude)
			//fmt.Printf("%f->%f\n", last, Latitude)
			//last = Longitude
			Longitude = randomBuffer(Longitude)
			//fmt.Printf("%f->%f\n", last, Longitude)
			_ = outputFileCsv.Exec("0", "0", convert.FloatToString(Longitude), convert.FloatToString(Latitude))
		}
	}
	_ = rawFile.Close()
	_ = outputFileCsv.Close()
}

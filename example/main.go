package main

import (
	"encoding/csv"
	"github.com/linger1216/go-gis/geom"
	"github.com/linger1216/go-gis/track"
	"github.com/linger1216/go-gis/visualizer"
	"github.com/linger1216/go-utils/convert"
	"io"
	"os"
)

func main() {

	filename := "res/xy.csv"

	f, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer func() { _ = f.Close() }()
	reader := csv.NewReader(f)
	reader.LazyQuotes = true

	xys := make([]*geom.LngLat, 0)
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		if len(line) >= 0 {
			xys = append(xys, &geom.LngLat{
				Longitude: convert.StringToFloat(line[0]),
				Latitude:  convert.StringToFloat(line[1]),
			})
		}
	}

	ops := &track.SimplifyOption{
		Degree: 0.000006,
	}
	visualizer.DrawScatter(1800, 900, "scatter.html", 0, geom.NewLineString(xys...))
	simple := track.Simplify(ops, geom.NewLineString(xys...))
	visualizer.DrawScatter(1800, 900, "scatter_simple.html", ops.Degree, simple)
}

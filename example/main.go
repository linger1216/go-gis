package main

import (
	"encoding/csv"
	"fmt"
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

	degrees := make([]float64, 0)
	for i := 1; i <= 5; i++ {
		degrees = append(degrees, float64(i))
	}
	for i := range degrees {
		ops := &track.SimplifyOption{
			Degree: degrees[i],
		}
		visualizer.DrawScatter(1200, 700, "scatter.html",
			fmt.Sprintf("visualizer %d points", len(xys)), xys...)
		simple := track.Simplify(ops, xys...)

		visualizer.DrawScatter(1200, 700, "scatter_simple.html",
			fmt.Sprintf("visualizer %d points with epsilon:%f", len(simple), ops.Degree), simple...)
		fmt.Printf("degress:%f origin:%d current:%d\n", degrees[i], len(xys), len(simple))
	}
}

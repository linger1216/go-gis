package main

import (
	"encoding/csv"
	"fmt"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/linger1216/go-gis/geom"
	"github.com/linger1216/go-gis/hub"
	"github.com/linger1216/go-gis/visualizer"
	"github.com/linger1216/go-utils/convert"
	"io"
	"os"
)

func main() {
	count := 0
	jump := 0
	filename := "/Users/lid/Downloads/cpp_proc/kalman-filter/res.csv"
	f, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	raw := make([]*geom.LngLat, 0)
	handled := make([]*geom.LngLat, 0)

	reader := csv.NewReader(f)
	reader.LazyQuotes = true
	for {
		line, err := reader.Read()
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
			raw = append(raw, &geom.LngLat{
				Latitude:  convert.StringToFloat(line[0]),
				Longitude: convert.StringToFloat(line[1]),
			})
			handled = append(handled, &geom.LngLat{
				Latitude:  convert.StringToFloat(line[2]),
				Longitude: convert.StringToFloat(line[3]),
			})
		}
	}

	_ = f.Close()

	rawPoints := make([]hub.TrackPointer, len(raw))
	for i := range raw {
		rawPoints[i] = raw[i]
	}

	handledPoints := make([]hub.TrackPointer, len(handled))
	for i := range handled {
		handledPoints[i] = handled[i]
	}

	width := 600
	height := 480

	page := components.NewPage()
	page.SetLayout(components.PageFlexLayout)

	// draw raw
	before := visualizer.DrawLine(width, height, fmt.Sprintf("visualizer %d rawPoints", len(rawPoints)), rawPoints...)
	page.AddCharts(before)

	// draw handle
	after := visualizer.DrawLine(width, height, fmt.Sprintf("visualizer %d handledPoints", len(handledPoints)), handledPoints...)
	page.AddCharts(after)

	pageFile, err := os.Create("kalman_page.html")
	if err != nil {
		panic(err)
	}
	_ = page.Render(pageFile)
}

package main

import (
	"encoding/csv"
	"fmt"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/linger1216/go-gis/model/geom"
	"github.com/linger1216/go-gis/model/hub"
	"github.com/linger1216/go-gis/src/track"
	"github.com/linger1216/go-gis/src/visualizer"
	"github.com/linger1216/go-utils/convert"
	"io"
	"os"
)

func main() {

	// common
	count := 0
	jump := 1
	visualizerFilename := "kalman-filter.html"

	// draw output
	width := 600
	height := 400

	page := components.NewPage()
	page.PageTitle = visualizerFilename
	page.SetLayout(components.PageFlexLayout)

	// raw
	rawFilename := "res/test2_line_RaPts1000_random.csv"
	rawLatIndex := 3
	rawLngIndex := 2
	rawFile, err := os.OpenFile(rawFilename, os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	rawXYs := make([]*geom.LngLat, 0)
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
			rawXYs = append(rawXYs, &geom.LngLat{
				Latitude:  convert.StringToFloat(line[rawLatIndex]),
				Longitude: convert.StringToFloat(line[rawLngIndex]),
			})
		}
	}
	_ = rawFile.Close()

	rawPoints := make([]hub.TrackPointer, len(rawXYs))
	for i := range rawXYs {
		rawPoints[i] = rawXYs[i]
	}

	rawPoints = rawPoints[:100]
	rawView := visualizer.DrawLine(width, height, fmt.Sprintf("visualizer %d rawPoints", len(rawXYs)), rawPoints)
	page.AddCharts(rawView)
	_ = rawView

	kf := track.NewDenoise()

	for i := 0; i <= 7; i++ {
		kfPoints := kf.Exec(&track.DenoiseOption{Degree: float64(i)}, rawPoints...)
		kfView := visualizer.DrawLine(width, height, fmt.Sprintf("go kf visualizer %d rawPoints with degress:%d", len(kfPoints), i), kfPoints)
		page.AddCharts(kfView)
		_ = kfView
	}

	pageFile, err := os.Create(visualizerFilename)
	if err != nil {
		panic(err)
	}
	_ = page.Render(pageFile)
}

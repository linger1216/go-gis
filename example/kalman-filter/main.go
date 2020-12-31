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
	jump := 0
	visualizerFilename := "kalman-filter.html"

	// control
	rawFilename := "res/2000.csv"
	rawLatIndex := 0
	rawLngIndex := 1
	cppFilename := "res/kfcpp.csv"
	cppLatIndex := 2
	cppLngIndex := 3

	// load rawView
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

	// load cpp
	cppFile, err := os.OpenFile(cppFilename, os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	cppXYs := make([]*geom.LngLat, 0)
	cppReader := csv.NewReader(cppFile)
	cppReader.LazyQuotes = true
	for {
		line, err := cppReader.Read()
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
			cppXYs = append(cppXYs, &geom.LngLat{
				Latitude:  convert.StringToFloat(line[cppLatIndex]),
				Longitude: convert.StringToFloat(line[cppLngIndex]),
			})
		}
	}
	_ = cppFile.Close()

	cppPoints := make([]hub.TrackPointer, len(cppXYs))
	for i := range cppXYs {
		cppPoints[i] = cppXYs[i]
	}

	// kf 1e-5
	kf := track.NewDenoise()
	kfPoints := kf.Exec(&track.DenoiseOption{Degree: 1}, rawPoints...)

	// draw output
	width := 600
	height := 400

	page := components.NewPage()
	page.PageTitle = visualizerFilename
	page.SetLayout(components.PageFlexLayout)
	rawView := visualizer.DrawLine(width, height, fmt.Sprintf("visualizer %d rawPoints", len(rawXYs)), rawPoints...)
	page.AddCharts(rawView)
	kfView := visualizer.DrawLine(width, height, fmt.Sprintf("go kf visualizer %d rawPoints", len(kfPoints)), kfPoints...)
	page.AddCharts(kfView)
	cppView := visualizer.DrawLine(width, height, fmt.Sprintf("cpp kf visualizer %d rawPoints", len(cppPoints)), cppPoints...)
	page.AddCharts(cppView)

	pageFile, err := os.Create(visualizerFilename)
	if err != nil {
		panic(err)
	}
	_ = page.Render(pageFile)
}

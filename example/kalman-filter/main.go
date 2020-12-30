package main

import (
	"encoding/csv"
	"fmt"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/linger1216/go-gis/algo/track"
	"github.com/linger1216/go-gis/geom"
	"github.com/linger1216/go-gis/hub"
	"github.com/linger1216/go-gis/visualizer"
	"github.com/linger1216/go-utils/convert"
	"gonum.org/v1/gonum/mat"
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

	// kf
	kf := track.NewKalManFilter(4, 2, 0)
	if len(rawPoints) > 0 {
		kf.StatePost = mat.NewDense(4, 1, []float64{rawPoints[0].Point().Latitude, rawPoints[0].Point().Longitude, 0, 0})
	}
	kf.TransitionMatrix = mat.NewDense(4, 4, []float64{
		1, 0, 1, 0,
		0, 1, 0, 1,
		0, 0, 1, 0,
		0, 0, 0, 1,
	})
	kf.MeasurementMatrix = mat.NewDense(2, 4, []float64{
		1, 0, 0, 0,
		0, 1, 0, 0,
	})
	// 协方差矩阵
	kf.ProcessNoiseCov = mat.NewDiagonalRect(4, 4, track.MakeMatValue(4, 1, 1e-5))
	kf.MeasurementNoiseCov = mat.NewDiagonalRect(2, 2, track.MakeMatValue(2, 1, 1e-5))
	kf.ErrorCovPost = mat.NewDense(4, 4, []float64{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	})

	kfPoints := make([]hub.TrackPointer, 0)
	for i := 0; i < len(rawPoints); i++ {
		prediction := kf.Predict(nil)
		p := geom.NewLngLat(prediction.At(1, 0), prediction.At(0, 0))
		kfPoints = append(kfPoints, p)
		measurement := mat.NewDense(2, 1, []float64{rawPoints[i].Point().Latitude, rawPoints[i].Point().Longitude})
		kf.Correct(measurement)
	}

	// draw output
	width := 600
	height := 400

	page := components.NewPage()
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

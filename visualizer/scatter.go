package visualizer

import (
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/linger1216/go-gis/geom"
	"math"
	"os"
)

func generateScatterItems(lineString *geom.LineString) []opts.ScatterData {
	data := make([]opts.ScatterData, len(lineString.Coordinates))
	for i := range lineString.Coordinates {
		data[i] = opts.ScatterData{
			Name:       fmt.Sprintf("%f,%f", lineString.Coordinates[i].Longitude, lineString.Coordinates[i].Latitude),
			Value:      []float64{lineString.Coordinates[i].Latitude, lineString.Coordinates[i].Longitude},
			SymbolSize: 20,
		}
	}
	return data
}

func DrawScatter(width, height int, filename string, epsilon float64, lineString *geom.LineString) {

	if width == 0 {
		width = 1800
	}

	if height == 0 {
		height = 900
	}

	if len(filename) == 0 {
		filename = "scatter.html"
	}

	lngMin := math.MaxFloat64
	lngMax := float64(0)
	latMin := math.MaxFloat64
	latMax := float64(0)

	for i := range lineString.Coordinates {

		if lineString.Coordinates[i].Latitude > latMax {
			latMax = lineString.Coordinates[i].Latitude
		}

		if lineString.Coordinates[i].Latitude < latMin {
			latMin = lineString.Coordinates[i].Latitude
		}

		if lineString.Coordinates[i].Longitude > lngMax {
			lngMax = lineString.Coordinates[i].Longitude
		}

		if lineString.Coordinates[i].Longitude < lngMin {
			lngMin = lineString.Coordinates[i].Longitude
		}
	}

	scatter := charts.NewScatter()
	scatter.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Width:  fmt.Sprintf("%dpx", width),
			Height: fmt.Sprintf("%dpx", height),
		}),
		charts.WithTitleOpts(opts.Title{
			Title: fmt.Sprintf("xy visualizer %d points with epsilon:%f", len(lineString.Coordinates), epsilon),
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show: true,
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Name: "Latitude",
			Type: "value",
			Min:  latMin,
			Max:  latMax,
			SplitLine: &opts.SplitLine{
				Show: true,
			},
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: "Longitude",
			Type: "value",
			Min:  lngMin,
			Max:  lngMax,
			SplitLine: &opts.SplitLine{
				Show: true,
			},
		}),

		charts.WithDataZoomOpts(
			//opts.DataZoom{
			//	Type:       "slider",
			//	XAxisIndex: 0,
			//	Start:      float32(latMin) / float32(latMax) * 100,
			//	End:        100,
			//},
			//opts.DataZoom{
			//	Type:       "slider",
			//	YAxisIndex: 0,
			//	Start:      float32(lngMin) / float32(lngMax) * 100,
			//	End:        100,
			//},
			opts.DataZoom{
				Type:       "inside",
				XAxisIndex: 0,
				Start:      float32(latMin) / float32(latMax) * 100,
				End:        100,
			},
			opts.DataZoom{
				Type:       "inside",
				YAxisIndex: 0,
				Start:      float32(lngMin) / float32(lngMax) * 100,
				End:        100,
			},
		),
	)
	scatter.AddSeries("xy", generateScatterItems(lineString))
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	_ = scatter.Render(f)
}

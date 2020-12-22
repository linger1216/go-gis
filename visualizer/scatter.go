package visualizer

import (
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/linger1216/go-gis/geom"
	"os"
)

func generateScatterItems(coords ...*geom.LngLat) []opts.ScatterData {
	data := make([]opts.ScatterData, len(coords))
	for i := range coords {
		data[i] = opts.ScatterData{
			Name:       fmt.Sprintf("%f,%f", coords[i].Longitude, coords[i].Latitude),
			Value:      []float64{coords[i].Latitude, coords[i].Longitude},
			SymbolSize: 20,
		}
	}
	return data
}

func DrawScatter(width, height int, outputFilename, title string, coords ...*geom.LngLat) {

	if width == 0 {
		width = 1800
	}

	if height == 0 {
		height = 900
	}

	if len(outputFilename) == 0 {
		outputFilename = "scatter.html"
	}

	box := geom.BoundingRect(coords...)
	lngMin := box.LeftBottom.Longitude
	lngMax := box.RightTop.Longitude
	latMin := box.LeftBottom.Latitude
	latMax := box.RightTop.Latitude

	scatter := charts.NewScatter()
	scatter.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Width:  fmt.Sprintf("%dpx", width),
			Height: fmt.Sprintf("%dpx", height),
		}),
		charts.WithTitleOpts(opts.Title{
			Title: title,
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
	scatter.AddSeries("xy", generateScatterItems(coords...))
	f, err := os.Create(outputFilename)
	if err != nil {
		panic(err)
	}
	_ = scatter.Render(f)
}

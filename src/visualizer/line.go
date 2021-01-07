package visualizer

import (
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/linger1216/go-gis/src/track"
)

func generateLineItems(coords ...track.TrackPointer) []opts.LineData {
	data := make([]opts.LineData, len(coords))
	for i := range coords {
		data[i] = opts.LineData{
			Name:       fmt.Sprintf("%d", i),
			Value:      []float64{coords[i].Point().Latitude, coords[i].Point().Longitude},
			SymbolSize: 10,
		}
	}
	return data
}

func DrawLine(width, height int, title string, coords ...[]track.TrackPointer) components.Charter {

	if width == 0 {
		width = 1800
	}

	if height == 0 {
		height = 900
	}

	box := track.BoundingRect(coords...)
	lngMin := box.LeftBottom.Longitude
	lngMax := box.RightTop.Longitude
	latMin := box.LeftBottom.Latitude
	latMax := box.RightTop.Latitude

	_ = lngMax
	_ = lngMin
	_ = latMax
	_ = latMin

	Line := charts.NewLine()
	Line.SetGlobalOptions(
		charts.WithColorsOpts(opts.Colors{"green", "red"}),
		charts.WithInitializationOpts(opts.Initialization{
			Width:  fmt.Sprintf("%dpx", width),
			Height: fmt.Sprintf("%dpx", height),
		}),
		charts.WithTitleOpts(opts.Title{
			Title: title,
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show:      true,
			Formatter: "[{b}]:{c}",
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

	for i := range coords {
		Line.AddSeries(fmt.Sprintf("cpp-%d", i), generateLineItems(coords[i]...))
	}

	Line.SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: true}))

	return Line
}

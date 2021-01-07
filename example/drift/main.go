package main

import (
	"encoding/csv"
	"fmt"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/linger1216/go-gis/src/geom"
	"github.com/linger1216/go-gis/src/track"
	"github.com/linger1216/go-gis/src/visualizer"
	"github.com/linger1216/go-utils/convert"
	"io"
	"os"
)

func main() {

	count := 0
	// control
	filename := "res/drift_random.csv"
	latIndex := 0
	lngIndex := 1
	jump := 0
	visualizerFilename := "drift.html"

	f, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	xys := make([]*geom.LngLat, 0)

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
			xys = append(xys, &geom.LngLat{
				Latitude:  convert.StringToFloat(line[latIndex]),
				Longitude: convert.StringToFloat(line[lngIndex]),
			})
		}
	}

	_ = f.Close()

	points := make([]track.TrackPointer, len(xys))
	for i := range xys {
		points[i] = xys[i]
	}

	width := 600
	height := 500

	page := components.NewPage()
	page.PageTitle = visualizerFilename
	page.SetLayout(components.PageFlexLayout)
	origin := visualizer.DrawLine(width, height, fmt.Sprintf("visualizer %d points", len(xys)), points)
	page.AddCharts(origin)

	d := track.NewDrift()
	simple := d.Exec(&track.DriftOption{SegmentPolicy: track.SegmentPolicyByInterval}, points...)
	page.AddCharts(visualizer.DrawLine(width, height, fmt.Sprintf("visualizer %d points after",
		len(simple)), simple))

	// 分段情况, 用于调试
	//arr := track.NewSegment().SegmentByDist(nil, points...)
	//for i := range arr {
	//	page.AddCharts(visualizer.DrawLine(width, height, fmt.Sprintf("visualizer %d points Segment[%d]-before", len(arr[i]), i), arr[i]))
	//	drifted := track.NewDrift().Exec(nil, arr[i]...)
	//	page.AddCharts(visualizer.DrawLine(width, height, fmt.Sprintf("visualizer %d points Segment[%d]-after", len(drifted), i), drifted))
	//}

	pageFile, err := os.Create(visualizerFilename)
	if err != nil {
		panic(err)
	}
	_ = page.Render(pageFile)
}

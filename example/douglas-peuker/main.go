package main

import (
	"bufio"
	"fmt"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/linger1216/go-gis/algo/track"
	"github.com/linger1216/go-gis/geom"
	"github.com/linger1216/go-gis/hub"
	"github.com/linger1216/go-gis/visualizer"
	"github.com/linger1216/go-utils/convert"
	"os"
	"strings"
)

func main() {

	count := 0
	jump := 6
	filename := "res/20080914080824.plt"
	f, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	r := bufio.NewReaderSize(f, 4096)
	xys := make([]*geom.LngLat, 0)

	for {
		l, e := r.ReadBytes('\n')
		if e != nil && len(l) == 0 {
			break
		}

		count++
		if count <= jump {
			continue
		}
		line := strings.Split(string(l), ",")
		if len(line) >= 0 {
			xys = append(xys, &geom.LngLat{
				Latitude:  convert.StringToFloat(line[0]),
				Longitude: convert.StringToFloat(line[1]),
			})
		}
	}
	_ = f.Close()

	points := make([]hub.TrackPointer, len(xys))
	for i := range xys {
		points[i] = xys[i]
	}

	degrees := make([]float64, 0)
	for i := 1; i <= 7; i++ {
		degrees = append(degrees, float64(i))
	}

	width := 600
	height := 400

	page := components.NewPage()
	page.SetLayout(components.PageFlexLayout)
	// draw
	origin := visualizer.DrawLine(width, height, fmt.Sprintf("visualizer %d points", len(xys)), points...)

	page.AddCharts(origin)

	for i := range degrees {
		ops := &track.SimplifyOption{
			Degree: degrees[i],
		}
		simple := track.Simplify(ops, points...)
		page.AddCharts(visualizer.DrawLine(width, height, fmt.Sprintf("visualizer %d points with epsilon:%f",
			len(simple), ops.Degree), simple...))
		fmt.Printf("degress:%f origin:%d current:%d\n", degrees[i], len(xys), len(simple))
	}
	pageFile, err := os.Create("page.html")
	if err != nil {
		panic(err)
	}
	_ = page.Render(pageFile)
}

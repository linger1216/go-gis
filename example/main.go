package main

import (
	"bufio"
	"fmt"
	"github.com/linger1216/go-gis/geom"
	"github.com/linger1216/go-gis/track"
	"github.com/linger1216/go-gis/visualizer"
	"github.com/linger1216/go-utils/convert"
	"os"
	"strings"
)

func main() {

	count := 0
	jump := 6
	filename := "/Users/lid/Downloads/Geolife Trajectories 1.3/Data/070/Trajectory/20081002003702.plt"
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

	degrees := make([]float64, 0)
	for i := 1; i <= 7; i++ {
		degrees = append(degrees, float64(i))
	}

	visualizer.DrawScatter(1200, 700, "scatter.html",
		fmt.Sprintf("visualizer %d points", len(xys)), xys...)

	for i := range degrees {
		ops := &track.SimplifyOption{
			Degree: degrees[i],
		}
		simple := track.Simplify(ops, xys...)
		visualizer.DrawScatter(1200, 700, fmt.Sprintf("scatter_simplify_%d.html", int(degrees[i])),
			fmt.Sprintf("visualizer %d points with epsilon:%f", len(simple), ops.Degree), simple...)
		fmt.Printf("degress:%f origin:%d current:%d\n", degrees[i], len(xys), len(simple))
	}
}

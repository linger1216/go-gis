package track

import (
	"github.com/linger1216/go-gis/src/geom"
	"math"
)

func BoundingRect(xys ...[]TrackPointer) *geom.Box {
	lngMin := math.MaxFloat64
	lngMax := float64(0)
	latMin := math.MaxFloat64
	latMax := float64(0)

	for _, coords := range xys {
		for i := range coords {
			if coords[i].Point().Latitude > latMax {
				latMax = coords[i].Point().Latitude
			}
			if coords[i].Point().Latitude < latMin {
				latMin = coords[i].Point().Latitude
			}
			if coords[i].Point().Longitude > lngMax {
				lngMax = coords[i].Point().Longitude
			}
			if coords[i].Point().Longitude < lngMin {
				lngMin = coords[i].Point().Longitude
			}
		}
	}

	return &geom.Box{
		LeftBottom: &geom.LngLat{
			Longitude: lngMin,
			Latitude:  latMin,
		},
		RightTop: &geom.LngLat{
			Longitude: lngMax,
			Latitude:  latMax,
		},
	}
}

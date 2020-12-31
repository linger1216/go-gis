package hub

import (
	"github.com/linger1216/go-gis/model/geom"
	"math"
)

func BoundingRect(coords ...TrackPointer) *geom.Box {
	lngMin := math.MaxFloat64
	lngMax := float64(0)
	latMin := math.MaxFloat64
	latMax := float64(0)

	for i := range coords {
		if coords[i].Position().Latitude > latMax {
			latMax = coords[i].Position().Latitude
		}
		if coords[i].Position().Latitude < latMin {
			latMin = coords[i].Position().Latitude
		}
		if coords[i].Position().Longitude > lngMax {
			lngMax = coords[i].Position().Longitude
		}
		if coords[i].Position().Longitude < lngMin {
			lngMin = coords[i].Position().Longitude
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

package geom

import (
	"github.com/linger1216/go-gis/hub"
	"math"
)

type Box struct {
	LeftBottom *LngLat `protobuf:"bytes,1,opt,name=left_bottom,json=leftBottom,proto3" json:"leftBottom,omitempty"`
	RightTop   *LngLat `protobuf:"bytes,2,opt,name=right_top,json=rightTop,proto3" json:"rightTop,omitempty"`
}

func BoundingRect(coords ...hub.TrackPointer) *Box {
	lngMin := math.MaxFloat64
	lngMax := float64(0)
	latMin := math.MaxFloat64
	latMax := float64(0)

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
	return &Box{
		LeftBottom: &LngLat{
			Longitude: lngMin,
			Latitude:  latMin,
		},
		RightTop: &LngLat{
			Longitude: lngMax,
			Latitude:  latMax,
		},
	}
}

package geom

import "math"

type Box struct {
	LeftBottom *LngLat `protobuf:"bytes,1,opt,name=left_bottom,json=leftBottom,proto3" json:"leftBottom,omitempty"`
	RightTop   *LngLat `protobuf:"bytes,2,opt,name=right_top,json=rightTop,proto3" json:"rightTop,omitempty"`
}

func BoundingRect(coords ...*LngLat) *Box {
	lngMin := math.MaxFloat64
	lngMax := float64(0)
	latMin := math.MaxFloat64
	latMax := float64(0)

	for i := range coords {
		if coords[i].Latitude > latMax {
			latMax = coords[i].Latitude
		}
		if coords[i].Latitude < latMin {
			latMin = coords[i].Latitude
		}
		if coords[i].Longitude > lngMax {
			lngMax = coords[i].Longitude
		}
		if coords[i].Longitude < lngMin {
			lngMin = coords[i].Longitude
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

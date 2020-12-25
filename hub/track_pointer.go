package hub

import (
	"github.com/linger1216/go-gis/geom"
	"math"
)

type TrackPointer interface {
	ID() string
	Point() *geom.LngLat
	Timestamp() int64
}

type TrackPoint struct {
	Id          string      `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty" bson:"_id"`
	Timestamp   int64       `protobuf:"varint,3,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	ReceiveTime int64       `protobuf:"varint,4,opt,name=receive_time,json=receiveTime,proto3" json:"receiveTime,omitempty"`
	AssetId     string      `protobuf:"bytes,1,opt,name=asset_id,proto3" json:"assetId,omitempty"`
	Point       *geom.Point `protobuf:"bytes,10,opt,name=point,proto3" json:"point,omitempty"`
}

func BoundingRect(coords ...TrackPointer) *geom.Box {
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

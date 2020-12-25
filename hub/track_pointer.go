package hub

import (
	"github.com/linger1216/go-gis/geom"
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

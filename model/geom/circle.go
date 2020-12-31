package geom

type Circle struct {
	Center *LngLat `protobuf:"bytes,1,opt,name=center,proto3" json:"center,omitempty"`
	Radius float32 `protobuf:"fixed32,2,opt,name=radius,proto3" json:"radius,omitempty"`
}

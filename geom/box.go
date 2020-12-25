package geom

type Box struct {
	LeftBottom *LngLat `protobuf:"bytes,1,opt,name=left_bottom,json=leftBottom,proto3" json:"leftBottom,omitempty"`
	RightTop   *LngLat `protobuf:"bytes,2,opt,name=right_top,json=rightTop,proto3" json:"rightTop,omitempty"`
}

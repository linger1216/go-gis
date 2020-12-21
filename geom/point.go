package geom

// Point
type Point struct {
	Coordinates      *LngLat `protobuf:"bytes,1,opt,name=coordinates,proto3" json:"coordinates,omitempty"`
	Type             string  `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
	SpatialReference string  `protobuf:"bytes,4,opt,name=spatial_reference,json=spatialReference,proto3" json:"spatialReference,omitempty"`
}

type MultiPoint struct {
	Points []*LngLat `protobuf:"bytes,1,rep,name=points,proto3" json:"points,omitempty"`
	Type   string    `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
}

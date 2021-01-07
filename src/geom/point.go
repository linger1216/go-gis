package geom

// Point
type Point struct {
	Coordinate       *LngLat `protobuf:"bytes,1,opt,name=coordinates,proto3" json:"coordinates,omitempty"`
	Type             string  `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
	SpatialReference string  `protobuf:"bytes,4,opt,name=spatial_reference,json=spatialReference,proto3" json:"spatialReference,omitempty"`
}

func NewPoint(spatialReference string, coordinates ...float64) *Point {
	return &Point{Coordinate: NewLngLat(coordinates...), Type: "Position", SpatialReference: spatialReference}
}

type MultiPoint struct {
	Points []*LngLat `protobuf:"bytes,1,rep,name=points,proto3" json:"points,omitempty"`
	Type   string    `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
}

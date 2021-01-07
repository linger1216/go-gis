package geom

// Polygon
type Polygon struct {
	LineStrings []*LineString `protobuf:"bytes,1,rep,name=line_strings,json=lineStrings,proto3" json:"lineStrings,omitempty"`
	Type        string        `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
}

type MultiPolygon struct {
	Polygons []*Polygon `protobuf:"bytes,1,rep,name=polygons,proto3" json:"polygons,omitempty"`
	Type     string     `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
}

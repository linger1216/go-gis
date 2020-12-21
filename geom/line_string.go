package geom

type LineString struct {
	Coordinates []*LngLat `protobuf:"bytes,1,rep,name=coordinates,proto3" json:"coordinates,omitempty"`
	Type        string    `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
}

type MultiLineString struct {
	LineStrings []*LineString `protobuf:"bytes,1,rep,name=line_strings,json=lineStrings,proto3" json:"lineStrings,omitempty"`
	Type        string        `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
}

func NewLineString(coordinates ...*LngLat) *LineString {
	line := &LineString{}
	line.Coordinates = append(line.Coordinates, coordinates...)
	return line
}

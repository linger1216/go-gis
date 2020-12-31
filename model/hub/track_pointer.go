package hub

import (
	"github.com/linger1216/go-gis/model/geom"
)

type TrackPointer interface {
	ID() string
	Point() *geom.LngLat
	Timestamp() int64
}

type TrackPoint struct {
	DeviceId string      `json:"deviceId,omitempty"`
	Time     int64       `json:"time,omitempty"`
	Position *geom.Point `json:"position,omitempty"`
}

func NewTrackPoint(deviceId string, time int64, lng, lat float64) *TrackPoint {
	return &TrackPoint{DeviceId: deviceId, Time: time, Position: geom.NewPoint(geom.Gcj02, lng, lat)}
}

func (t *TrackPoint) ID() string {
	return t.DeviceId
}

func (t *TrackPoint) Point() *geom.LngLat {
	return t.Position.Coordinate
}

func (t *TrackPoint) Timestamp() int64 {
	return t.Time
}

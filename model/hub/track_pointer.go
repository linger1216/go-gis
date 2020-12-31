package hub

import (
	"github.com/linger1216/go-gis/model/geom"
)

type TrackPointer interface {
	ID() string
	Position() *geom.LngLat
	Timestamp() int64
}

type TrackPoint struct {
	DeviceId string      `json:"deviceId,omitempty"`
	Time     int64       `json:"timestamp,omitempty"`
	Point    *geom.Point `json:"point,omitempty"`
}

func NewTrackPoint(deviceId string, time int64, lng, lat float64) *TrackPoint {
	return &TrackPoint{DeviceId: deviceId, Time: time, Point: geom.NewPoint(geom.Gcj02, lng, lat)}
}

func (t *TrackPoint) ID() string {
	return t.DeviceId
}

func (t *TrackPoint) Position() *geom.LngLat {
	return t.Point.Coordinate
}

func (t *TrackPoint) Timestamp() int64 {
	return t.Time
}

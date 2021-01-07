package track

import (
	"github.com/linger1216/go-gis/src/geom"
)

const (
	IntervalTimes   = 3 // 超过3倍常规时间就认为不是同一段轨迹
	DistanceTimes   = 3 // 超过3倍常规距离就认为不是同一段轨迹
	MinSegmentPoint = 5 // 最少的线段组成的点
)

type SegmentOption struct {
}

type Segment struct {
}

func NewSegment() *Segment {
	return &Segment{}
}

func (d *Segment) SegmentByInterval(ops *SegmentOption, coords ...TrackPointer) [][]TrackPointer {
	size := len(coords)
	ret := make([][]TrackPointer, 0)
	ret = append(ret, []TrackPointer{})
	lastInterval := int64(0)
	lastPos := 0
	i := lastPos + 1
	for i = lastPos + 1; i < size; i++ {
		interval := coords[i].Timestamp() - coords[i-1].Timestamp()
		if lastInterval > 0 && interval >= lastInterval*IntervalTimes {
			ret[len(ret)-1] = append(ret[len(ret)-1], coords[lastPos:i]...)
			lastInterval = 0
			lastPos = i
			ret = append(ret, []TrackPointer{})
			continue
		}
		lastInterval = interval
	}
	ret[len(ret)-1] = append(ret[len(ret)-1], coords[lastPos:i]...)
	if len(ret[len(ret)-1]) == 0 {
		ret = ret[:len(ret)-1]
	}
	return ret
}

func (d *Segment) SegmentByDist(ops *SegmentOption, coords ...TrackPointer) [][]TrackPointer {
	size := len(coords)
	ret := make([][]TrackPointer, 0)
	ret = append(ret, []TrackPointer{})
	lastDist := float64(0)
	lastPos := 0
	i := lastPos + 1
	for i = lastPos + 1; i < size; i++ {
		dist := geom.EuclideanDistance(coords[i].Point().Longitude, coords[i].Point().Latitude,
			coords[i-1].Point().Longitude, coords[i-1].Point().Latitude)
		if lastDist > 0 && dist >= lastDist*DistanceTimes && i-lastPos+1 >= MinSegmentPoint {
			ret[len(ret)-1] = append(ret[len(ret)-1], coords[lastPos:i+1]...)
			lastDist = 0
			lastPos = i + 1
			ret = append(ret, []TrackPointer{})
			continue
		}
		lastDist = dist
	}
	ret[len(ret)-1] = append(ret[len(ret)-1], coords[lastPos:size]...)

	if len(ret[len(ret)-1]) == 0 {
		ret = ret[:len(ret)-1]
	}

	return ret
}

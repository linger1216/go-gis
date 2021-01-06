package track

import (
	"github.com/linger1216/go-gis/model/geom"
	"github.com/linger1216/go-gis/model/hub"
	"math"
)

const (
	SegmentPolicyByDist     = "Dist"
	SegmentPolicyByInterval = "Interval"
	NoSegment               = "Nope"
)

type DriftOption struct {
	SegmentPolicy string // "Dist", "Interval", "Nope"
}

type Drift struct {
}

func NewDrift() *Drift {
	return &Drift{}
}

func (d *Drift) Exec(ops *DriftOption, coords ...hub.TrackPointer) []hub.TrackPointer {

	if ops == nil {
		ops = &DriftOption{NoSegment}
	}

	var arr [][]hub.TrackPointer

	switch ops.SegmentPolicy {
	case SegmentPolicyByDist:
		arr = NewSegment().SegmentByDist(nil, coords...)
	case SegmentPolicyByInterval:
		arr = NewSegment().SegmentByInterval(nil, coords...)
	default:
		arr = append(arr, coords)
	}

	ret := make([]hub.TrackPointer, 0, len(coords))
	for i := range arr {
		if v := d._DriftByStandardDeviation(ops, arr[i]...); len(v) > 0 {
			if len(v) == 2 {
				// todo
				// 暂时不知道要放弃哪个点
			}
			ret = append(ret, v...)
		}
	}
	return ret
}

func (d *Drift) _DriftByStandardDeviation(ops *DriftOption, coords ...hub.TrackPointer) []hub.TrackPointer {
	_ = ops
	size := len(coords)
	if size < 3 {
		return coords
	}

	dists := make([]float64, len(coords))
	sum := float64(0)
	for i := 1; i < size; i++ {
		dist := geom.EuclideanDistance(coords[i].Point().Longitude, coords[i].Point().Latitude,
			coords[i-1].Point().Longitude, coords[i-1].Point().Latitude)
		dists[i] = dist
		sum += dist
	}
	avg := sum / float64(size-1)

	E := float64(0)
	for i := 1; i < size; i++ {
		E += (dists[i] - avg) * (dists[i] - avg)
	}

	E /= float64(size)
	standardDeviation := math.Sqrt(E)
	_ = standardDeviation

	// 点数越少, 对标准差的影响越大
	// 找最值的方法来删除
	if size < 4 {
		// 方案1
		// 找到1个点到其他点距离最大, 那么就删除这个点, 不咋地
		//maxDist := float64(0)
		//maxIndex := 0
		//for i := 0; i < 3; i++ {
		//	dist := float64(0)
		//	for j := 0; j < 3; j++ {
		//		if i != j {
		//			dist += geom.EuclideanDistance(coords[i].Point().Longitude, coords[i].Point().Latitude,
		//				coords[j].Point().Longitude, coords[j].Point().Latitude)
		//		}
		//	}
		//	if dist > maxDist {
		//		maxDist = dist
		//		maxIndex = i
		//	}
		//}
		//ret := make([]hub.TrackPointer, 0, 3)
		//ret = append(ret, coords[0:maxIndex]...)
		//ret = append(ret, coords[maxIndex+1:]...)
		//return ret

		// 方案2
		// 找一个点和标准差的距离最大效果也不太好
		//maxIndex := size-1
		//maxValue := float64(0)
		//for i := 1; i < size; i++ {
		//	v := math.Abs(dists[i]-avg) - standardDeviation
		//	// 本质上>0都是不符合条件的
		//	if v > 0 && v > maxValue{
		//		maxValue = v
		//		maxIndex = i
		//	}
		//}
		//ret := make([]hub.TrackPointer, 0, 3)
		//ret = append(ret, coords[0:maxIndex]...)
		//ret = append(ret, coords[maxIndex+1:]...)
		//return ret

		// 方案3
		// 因为默认是发生异常点就分段, 所以大概率第1个点是标准点
		// 所以简单判断哪个距离最大是放弃哪个
		maxIndex := size - 1
		maxValue := float64(0)
		for i := 1; i < size; i++ {
			if dists[i] > maxValue {
				maxValue = dists[i]
				maxIndex = i
			}
		}
		ret := make([]hub.TrackPointer, 0, 3)
		ret = append(ret, coords[0:maxIndex]...)
		ret = append(ret, coords[maxIndex+1:]...)
		return ret
	}

	ret := make([]hub.TrackPointer, 0, size)
	ret = append(ret, coords[0])
	for i := 1; i < size; i++ {
		if dists[i]-avg <= standardDeviation {
			ret = append(ret, coords[i])
		}
	}
	return ret
}

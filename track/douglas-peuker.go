package track

import (
	"github.com/linger1216/go-gis/geom"
	"math"
)

type SimplifyOption struct {
	Degree float64 `protobuf:"fixed64,1,opt,name=degree,proto3" json:"degree,omitempty"`
}

// p到线段p1,p2的距离
// slope: 斜率
// intercept: 其中的b (y = kx + b)
// result: 距离
// 公式:
// https://zhuanlan.zhihu.com/p/26307123
// https://baike.baidu.com/item/%E4%B8%A4%E7%82%B9%E5%BC%8F

func findPerpendicularDistance(p, p1, p2 *geom.LngLat) (result float64) {
	if p1.Latitude == p2.Latitude {
		result = math.Abs(p.Latitude - p1.Latitude)
	} else {
		slope := (p2.Longitude - p1.Longitude) / (p2.Latitude - p1.Latitude)
		intercept := p1.Longitude - (slope * p1.Latitude)
		result = math.Abs(slope*p.Latitude-p.Longitude+intercept) / math.Sqrt(math.Pow(slope, 2)+1)
	}
	return
}

func Simplify(ops *SimplifyOption, line *geom.LineString) *geom.LineString {
	coords := _simplify(ops.Degree, line.Coordinates...)
	return &geom.LineString{
		Coordinates: coords,
	}
}

func _simplify(epsilon float64, coords ...*geom.LngLat) []*geom.LngLat {
	if len(coords) < 3 {
		return coords
	}
	firstPoint := coords[0]
	lastPoint := coords[len(coords)-1]
	index := -1
	dist := float64(0)
	for i := 1; i < len(coords)-1; i++ {
		cDist := findPerpendicularDistance(coords[i], firstPoint, lastPoint)
		if cDist > dist {
			dist = cDist
			index = i
		}
	}
	if dist > epsilon {
		l1 := coords[0 : index+1]
		l2 := coords[index:]
		r1 := _simplify(epsilon, l1...)
		r2 := _simplify(epsilon, l2...)
		rs := append(r1[0:len(r1)-1], r2...)
		return rs
	} else {
		ret := make([]*geom.LngLat, 0)
		ret = append(ret, firstPoint, lastPoint)
		return ret
	}
}

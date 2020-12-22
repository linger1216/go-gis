package track

import (
	"github.com/linger1216/go-gis/geom"
	"github.com/linger1216/go-utils/algorithm"
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
func _perpendicularDistance(p, p1, p2 *geom.LngLat) (result float64) {
	if p1.Latitude == p2.Latitude {
		result = math.Abs(p.Latitude - p1.Latitude)
	} else {
		slope := (p2.Longitude - p1.Longitude) / (p2.Latitude - p1.Latitude)
		intercept := p1.Longitude - (slope * p1.Latitude)
		result = math.Abs(slope*p.Latitude-p.Longitude+intercept) / math.Sqrt(math.Pow(slope, 2)+1)
	}
	return
}

/*
*
Simplify 使用dp算法进行抽稀
抽稀力度由SimplifyOption字段中的degree选项指定。
Degree 取值 [1,7]: 其中1代表抽稀力度最弱,7最强
以某一段轨迹抽稀取值后结果如下, 仅供参考:
degress:1 total:937 current:494
degress:2 total:937 current:352
degress:3 total:937 current:256
degress:4 total:937 current:148
degress:5 total:937 current:94
degress:6 total:937 current:53
degress:7 total:937 current:29
*
*/
func Simplify(ops *SimplifyOption, coords ...*geom.LngLat) []*geom.LngLat {
	epsilon := _transEpsilon(int(ops.Degree))
	return _douglasPeucker(epsilon, coords...)
}

func _transEpsilon(level int) float64 {
	switch level {
	case 1:
		return 0.000003
	case 2:
		return 0.000010
	case 3:
		return 0.000020
	case 4:
		return 0.000050
	case 5:
		return 0.000100
	case 6:
		return 0.000200
	case 7:
		return 0.000500
	default:
		return 0.000030
	}
}

func _douglasPeuckerRecursion(epsilon float64, coords ...*geom.LngLat) []*geom.LngLat {
	if len(coords) < 3 {
		return coords
	}
	firstPoint := coords[0]
	lastPoint := coords[len(coords)-1]
	index := -1
	dist := float64(0)
	for i := 1; i < len(coords)-1; i++ {
		cDist := _perpendicularDistance(coords[i], firstPoint, lastPoint)
		if cDist > dist {
			dist = cDist
			index = i
		}
	}
	if dist > epsilon {
		l1 := coords[0 : index+1]
		l2 := coords[index:]
		r1 := _douglasPeuckerRecursion(epsilon, l1...)
		r2 := _douglasPeuckerRecursion(epsilon, l2...)
		rs := append(r1[0:len(r1)-1], r2...)
		return rs
	} else {
		ret := make([]*geom.LngLat, 0)
		ret = append(ret, firstPoint, lastPoint)
		return ret
	}
}

func _douglasPeucker(epsilon float64, coords ...*geom.LngLat) []*geom.LngLat {
	if len(coords) < 3 {
		return coords
	}
	var size = len(coords)
	markers := make([]int, size)
	first := 0
	last := size - 1

	stack := algorithm.Stack(make([]interface{}, 0))
	i, index := 0, 0
	maxDistance, distance := float64(0), float64(0)
	markers[first], markers[last] = 1, 1
	stack.Push(first)
	stack.Push(last)

	for !stack.IsEmpty() {
		if v, err := stack.Pop(); err == nil {
			last = v.(int)
		} else {
			break
		}
		if v, err := stack.Pop(); err == nil {
			first = v.(int)
		} else {
			break
		}
		maxDistance = 0
		for i = first + 1; i < last; i++ {
			distance = _perpendicularDistance(coords[i], coords[first], coords[last])
			if distance > maxDistance {
				index = i
				maxDistance = distance
			}
		}

		if maxDistance > epsilon {
			markers[index] = 1
			stack.Push(first)
			stack.Push(index)
			stack.Push(index)
			stack.Push(last)
		}
	}
	ret := make([]*geom.LngLat, 0, size)
	for i = 0; i < size; i++ {
		if markers[i] > 0 {
			ret = append(ret, coords[i])
		}
	}
	return ret
}
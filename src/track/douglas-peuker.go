package track

import (
	"github.com/linger1216/go-utils/algorithm"
	"math"
)

// p到线段p1,p2的距离
// slope: 斜率
// intercept: 其中的b (y = kx + b)
// result: 距离
// 公式:
// https://zhuanlan.zhihu.com/p/26307123
// https://baike.baidu.com/item/%E4%B8%A4%E7%82%B9%E5%BC%8F
func _perpendicularDistance(p, p1, p2 TrackPointer) (result float64) {
	if p1.Point().Latitude == p2.Point().Latitude {
		result = math.Abs(p.Point().Latitude - p1.Point().Latitude)
	} else {
		slope := (p2.Point().Longitude - p1.Point().Longitude) / (p2.Point().Latitude - p1.Point().Latitude)
		intercept := p1.Point().Longitude - (slope * p1.Point().Latitude)
		result = math.Abs(slope*p.Point().Latitude-p.Point().Longitude+intercept) / math.Sqrt(math.Pow(slope, 2)+1)
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

func DouglasPeuckerRecursion(epsilon float64, coords ...TrackPointer) []TrackPointer {
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
		r1 := DouglasPeuckerRecursion(epsilon, l1...)
		r2 := DouglasPeuckerRecursion(epsilon, l2...)
		rs := append(r1[0:len(r1)-1], r2...)
		return rs
	} else {
		ret := make([]TrackPointer, 0)
		ret = append(ret, firstPoint, lastPoint)
		return ret
	}
}

func DouglasPeucker(epsilon float64, coords ...TrackPointer) []TrackPointer {
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
	ret := make([]TrackPointer, 0, size)
	for i = 0; i < size; i++ {
		if markers[i] > 0 {
			ret = append(ret, coords[i])
		}
	}
	return ret
}

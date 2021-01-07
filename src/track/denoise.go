package track

import (
	"fmt"
	"github.com/linger1216/go-gis/src/algo"
	"github.com/linger1216/go-gis/src/geom"
	"gonum.org/v1/gonum/mat"
	"math"
)

const (
	MeasurementNoise = 2e-5
)

type DenoiseOption struct {
	Degree float64 `protobuf:"fixed64,1,opt,name=degree,proto3" json:"degree,omitempty"`
}

type NormalDenoise struct {
	kf *algo.KalManFilter
}

// option *DenoiseOption
func NewNormalDenoise() *NormalDenoise {
	// kf
	kf := algo.NewKalManFilter(4, 2, 0)

	kf.TransitionMatrix = mat.NewDense(4, 4, []float64{
		1, 0, 1, 0,
		0, 1, 0, 1,
		0, 0, 1, 0,
		0, 0, 0, 1,
	})
	kf.MeasurementMatrix = mat.NewDense(2, 4, []float64{
		1, 0, 0, 0,
		0, 1, 0, 0,
	})

	kf.ErrorCovPost = mat.NewDense(4, 4, []float64{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	})

	return &NormalDenoise{kf: kf}
}

func (d *NormalDenoise) Exec(ops *DenoiseOption, coords ...TrackPointer) []TrackPointer {
	ret := make([]TrackPointer, 0, len(coords))
	tracks := d._part(coords...)
	for i := range tracks {
		if points := d._predict(ops, tracks[i]...); len(points) > 0 {
			ret = append(ret, points...)
		}
	}
	return ret
}

func (d *NormalDenoise) _part(coords ...TrackPointer) [][]TrackPointer {

	m := make(map[int64]int64)
	lastTime := int64(0)
	for i := range coords {
		if lastTime == 0 {
			lastTime = coords[i].Timestamp()
		} else {
			dist := coords[i].Timestamp() - lastTime
			m[dist]++
			lastTime = coords[i].Timestamp()
		}
	}

	// has no timestamp
	if len(m) == 0 {
		return [][]TrackPointer{coords}
	}

	maxKey := int64(0)
	maxVal := int64(0)

	for k, v := range m {
		if v > maxVal {
			maxKey = k
			maxVal = v
		}
	}

	lastTime = 0
	maxDist := maxKey
	ret := make([][]TrackPointer, 0)
	ret = append(ret, []TrackPointer{})
	for i := range coords {
		if lastTime == 0 {
			lastTime = coords[i].Timestamp()
		} else {
			dist := coords[i].Timestamp() - lastTime
			if dist > maxDist {
				ret = append(ret, []TrackPointer{})
			}
			ret[len(ret)-1] = append(ret[len(ret)-1], coords[i])
			lastTime = coords[i].Timestamp()
		}
	}
	return ret
}

func (d *NormalDenoise) EpsilonString(level int) string {
	Q, R := d._transEpsilon(level)
	return fmt.Sprintf("Q:%.2fm, R:%.2fm", Q*1e6, R*1e6)
}

func (d *NormalDenoise) _transEpsilon(level int) (float64, float64) {
	switch level {
	case 1:
		return 3e-6, MeasurementNoise
	case 2:
		return 1e-6, MeasurementNoise
	case 3:
		return 5e-7, MeasurementNoise
	case 4:
		return 3e-7, MeasurementNoise
	case 5:
		return 1e-7, MeasurementNoise
	case 6:
		return 1e-8, MeasurementNoise
	case 7:
		return 1e-9, MeasurementNoise
	default:
		return -1, MeasurementNoise
	}
}

func (d *NormalDenoise) _predict(ops *DenoiseOption, coords ...TrackPointer) []TrackPointer {
	Q, R := d._transEpsilon(int(ops.Degree))
	if Q < 0 {
		//d.kf.ProcessNoiseCov = mat.NewDiagonalRect(4, 4, algo.MakeMatValue(4, 1, 0))
		//points := d.__predict(coords...)
		points := coords
		// todo
		// dists 没有方向方面的矢量运算
		dists := make([]float64, len(points))
		sum := float64(0)
		size := len(points)
		for i := range points {
			if i == 0 {
				dists[i] = 0
				continue
			}
			dist := geom.Distance(points[i].Point().Longitude, points[i].Point().Latitude,
				points[i-1].Point().Longitude, points[i-1].Point().Latitude)
			dists[i] = dist
			sum += dist
		}

		avg := sum / float64(size-1)
		E := float64(0)
		for i := 1; i < size; i++ {
			fmt.Printf("dist[%d]=%f\n", i, dists[i])
			E += (avg - dists[i]) * (avg - dists[i])
		}
		E /= float64(size) * 1e6
		Q = math.Abs(E - R)
	}
	d.kf.ProcessNoiseCov = mat.NewDiagonalRect(4, 4, algo.MakeMatValue(4, 1, Q))
	d.kf.MeasurementNoiseCov = mat.NewDiagonalRect(2, 2, algo.MakeMatValue(2, 1, R))
	return d.__predict(coords...)
}

func (d *NormalDenoise) __predict(coords ...TrackPointer) []TrackPointer {
	if len(coords) > 0 {
		d.kf.StatePost = mat.NewDense(4, 1, []float64{coords[0].Point().Latitude, coords[0].Point().Longitude, 0, 0})
	}
	kfPoints := make([]TrackPointer, 0)
	for i := 0; i < len(coords); i++ {
		prediction := d.kf.Predict(nil)
		p := geom.NewLngLat(prediction.At(1, 0), prediction.At(0, 0))
		kfPoints = append(kfPoints, p)
		measurement := mat.NewDense(2, 1, []float64{coords[i].Point().Latitude, coords[i].Point().Longitude})
		d.kf.Correct(measurement)
	}
	return kfPoints
}

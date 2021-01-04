package track

import (
	"fmt"
	"github.com/linger1216/go-gis/model/geom"
	"github.com/linger1216/go-gis/model/hub"
	"github.com/linger1216/go-gis/src/algo"
	"gonum.org/v1/gonum/mat"
)

type DenoiseOption struct {
	Degree float64 `protobuf:"fixed64,1,opt,name=degree,proto3" json:"degree,omitempty"`
}

type Denoise struct {
	kf *algo.KalManFilter
}

// option *DenoiseOption
func NewDenoise() *Denoise {
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

	return &Denoise{kf: kf}
}

func (d *Denoise) Exec(ops *DenoiseOption, coords ...hub.TrackPointer) []hub.TrackPointer {
	ret := make([]hub.TrackPointer, 0, len(coords))
	tracks := d._part(coords...)
	for i := range tracks {
		if points := d._predict(ops, tracks[i]...); len(points) > 0 {
			ret = append(ret, points...)
		}
	}
	return ret
}

func (d *Denoise) _part(coords ...hub.TrackPointer) [][]hub.TrackPointer {

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
		return [][]hub.TrackPointer{coords}
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
	ret := make([][]hub.TrackPointer, 0)
	ret = append(ret, []hub.TrackPointer{})
	for i := range coords {
		if lastTime == 0 {
			lastTime = coords[i].Timestamp()
		} else {
			dist := coords[i].Timestamp() - lastTime
			if dist > maxDist {
				ret = append(ret, []hub.TrackPointer{})
			}
			ret[len(ret)-1] = append(ret[len(ret)-1], coords[i])
			lastTime = coords[i].Timestamp()
		}
	}
	return ret
}

func (d *Denoise) EpsilonString(level int) string {
	Q, R := d._transEpsilon(level)
	return fmt.Sprintf("Q:%.2fm, R:%.2fm", Q*1e6, R*1e6)
}

func (d *Denoise) _transEpsilon(level int) (float64, float64) {
	R := 2e-5
	switch level {
	case 1:
		return 3e-6, R
	case 2:
		return 1e-6, R
	case 3:
		return 5e-7, R
	case 4:
		return 3e-7, R
	case 5:
		return 1e-7, R
	case 6:
		return 1e-8, R
	case 7:
		return 1e-9, R
	default:
		return 3e-7, R
	}
}

func (d *Denoise) _predict(ops *DenoiseOption, coords ...hub.TrackPointer) []hub.TrackPointer {
	Q, R := d._transEpsilon(int(ops.Degree))
	d.kf.ProcessNoiseCov = mat.NewDiagonalRect(4, 4, algo.MakeMatValue(4, 1, Q))
	d.kf.MeasurementNoiseCov = mat.NewDiagonalRect(2, 2, algo.MakeMatValue(2, 1, R))

	if len(coords) > 0 {
		d.kf.StatePost = mat.NewDense(4, 1, []float64{coords[0].Point().Latitude, coords[0].Point().Longitude, 0, 0})
	}
	kfPoints := make([]hub.TrackPointer, 0)
	for i := 0; i < len(coords); i++ {
		prediction := d.kf.Predict(nil)
		p := geom.NewLngLat(prediction.At(1, 0), prediction.At(0, 0))
		kfPoints = append(kfPoints, p)
		measurement := mat.NewDense(2, 1, []float64{coords[i].Point().Latitude, coords[i].Point().Longitude})
		d.kf.Correct(measurement)
	}
	return kfPoints
}

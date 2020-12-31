package algo

import (
	"fmt"
	"gonum.org/v1/gonum/mat"
)

type KalManFilter struct {
	statePre            *mat.Dense     //!< predicted state (x'(k)): x(k)=A*x(k-1)+B*u(k)
	StatePost           *mat.Dense     //!< corrected state (x(k)): x(k)=x'(k)+K(k)*(z(k)-H*x'(k))
	TransitionMatrix    *mat.Dense     //!< state transition matrix (A)
	ControlMatrix       *mat.Dense     //!< control matrix (B) (not used if there is no control)
	MeasurementMatrix   *mat.Dense     //!< measurement matrix (H)
	ProcessNoiseCov     *mat.BandDense //!< process noise covariance matrix (Q)
	MeasurementNoiseCov *mat.BandDense //!< measurement noise covariance matrix (R)
	ErrorCovPre         *mat.Dense     //!< priori error estimate covariance matrix (P'(k)): P'(k)=A*P(k-1)*At + Q)*/
	gain                *mat.Dense     //!< Kalman gain matrix (K(k)): K(k)=P'(k)*Ht*inv(H*P'(k)*Ht+R)
	ErrorCovPost        *mat.Dense     //!< posteriori error estimate covariance matrix (P(k)): P(k)=(I-K(k)*H)*P'(k)

	temp1 *mat.Dense
	temp2 *mat.Dense
	temp3 *mat.Dense
	temp4 *mat.Dense
	temp5 *mat.Dense
}

func MakeMatValue(row, col int, v float64) []float64 {
	size := row * col
	ret := make([]float64, size)
	for i := 0; i < size; i++ {
		ret[i] = v
	}
	return ret
}

func NewKalManFilter(dynamParams, measureParams, controlParams int) *KalManFilter {

	if dynamParams <= 0 || measureParams <= 0 {
		panic("todo error 1")
	}

	ret := &KalManFilter{}
	ret.statePre = mat.NewDense(dynamParams, 1, MakeMatValue(dynamParams, 1, 0))
	ret.StatePost = mat.NewDense(dynamParams, 1, MakeMatValue(dynamParams, 1, 0))
	ret.TransitionMatrix = mat.NewDense(dynamParams, dynamParams, MakeMatValue(
		dynamParams, dynamParams, 0))
	ret.ProcessNoiseCov = mat.NewDiagonalRect(dynamParams, dynamParams, MakeMatValue(dynamParams, 1, 1))
	ret.MeasurementMatrix = mat.NewDense(measureParams, dynamParams, MakeMatValue(measureParams, dynamParams, 0))
	ret.MeasurementNoiseCov = mat.NewDiagonalRect(measureParams, measureParams, MakeMatValue(measureParams, 1, 1))

	ret.ErrorCovPre = mat.NewDense(dynamParams, dynamParams, MakeMatValue(dynamParams, dynamParams, 0))
	ret.ErrorCovPost = mat.NewDense(dynamParams, dynamParams, MakeMatValue(dynamParams, dynamParams, 0))

	ret.gain = mat.NewDense(dynamParams, measureParams, MakeMatValue(dynamParams, measureParams, 0))

	if controlParams > 0 {
		ret.ControlMatrix = mat.NewDense(dynamParams, controlParams, MakeMatValue(dynamParams, controlParams, 0))
	}

	ret.temp1 = mat.NewDense(dynamParams, dynamParams, MakeMatValue(dynamParams, dynamParams, 0))
	ret.temp2 = mat.NewDense(measureParams, dynamParams, MakeMatValue(measureParams, dynamParams, 0))
	ret.temp3 = mat.NewDense(measureParams, measureParams, MakeMatValue(measureParams, measureParams, 0))
	ret.temp4 = mat.NewDense(measureParams, dynamParams, MakeMatValue(measureParams, dynamParams, 0))
	ret.temp5 = mat.NewDense(measureParams, 1, MakeMatValue(measureParams, 1, 0))

	return ret
}

func p(title string, dense mat.Matrix) {

	defer func() { fmt.Printf("\n") }()
	fmt.Printf("%s:\n", title)
	switch x := dense.(type) {
	case *mat.Dense:
		if x == nil {
			return
		}
	case *mat.BandDense:
		if x == nil {
			return
		}
	}
	row, col := dense.Dims()
	for i := 0; i < row; i++ {
		for j := 0; j < col; j++ {
			fmt.Printf("%f ", dense.At(i, j))
		}
		fmt.Printf("\n")
	}
}

func (k *KalManFilter) P() {
	p("statePre", k.statePre)
	p("StatePost", k.StatePost)
	p("TransitionMatrix", k.TransitionMatrix)
	p("ControlMatrix", k.ControlMatrix)
	p("MeasurementMatrix", k.MeasurementMatrix)
	p("ProcessNoiseCov", k.ProcessNoiseCov)
	p("MeasurementNoiseCov", k.MeasurementNoiseCov)
	p("ErrorCovPre", k.ErrorCovPre)
	p("gain", k.gain)
	p("ErrorCovPost", k.ErrorCovPost)

	p("temp1", k.temp1)
	p("temp2", k.temp2)
	p("temp3", k.temp3)
	p("temp4", k.temp4)
	p("temp5", k.temp5)
}

// 控制变量
func (k *KalManFilter) Predict(control *mat.Dense) *mat.Dense {
	// update the state: x'(k) = A*x(k)
	k.statePre = _mul(k.TransitionMatrix, k.StatePost)
	if control != nil {
		// x'(k) = x'(k) + B*u(k)
		k.StatePost.Add(k.StatePost, _mul(k.ControlMatrix, control))
	}
	// update error covariance matrices: temp1 = A*P(k)
	k.temp1 = _mul(k.TransitionMatrix, k.ErrorCovPost)
	// P'(k) = temp1*At + Q
	k.ErrorCovPre = _add(_mul(k.temp1, k.TransitionMatrix.T()), k.ProcessNoiseCov)
	k.StatePost.CloneFrom(k.statePre)
	k.ErrorCovPost.CloneFrom(k.ErrorCovPre)

	//k.P()
	return k.statePre
}

func (k *KalManFilter) Correct(measurement *mat.Dense) *mat.Dense {
	// temp2 = H*P'(k)
	k.temp2 = _mul(k.MeasurementMatrix, k.ErrorCovPre)
	// temp3 = temp2*Ht + R
	k.temp3 = _add(_mul(k.temp2, k.MeasurementMatrix.T()), k.MeasurementNoiseCov)
	// temp4 = inv(temp3)*temp2 = Kt(k)
	k.temp4 = _mul(_inverse(k.temp3), k.temp2)
	// K(k)
	k.gain = mat.DenseCopyOf(k.temp4.T())
	// temp5 = z(k) - H*x'(k)
	k.temp5 = _sub(measurement, _mul(k.MeasurementMatrix, k.statePre))
	// x(k) = x'(k) + K(k)*temp5
	k.StatePost = _add(k.statePre, _mul(k.gain, k.temp5))
	// P(k) = P'(k) - K(k)*temp2
	k.ErrorCovPost = _sub(k.ErrorCovPre, _mul(k.gain, k.temp2))
	//k.P()
	return k.StatePost
}

func _inverse(a *mat.Dense) *mat.Dense {
	ret := &mat.Dense{}
	err := ret.Inverse(a)
	if err != nil {
		panic(err)
	}
	return ret
}

func _mul(a, b mat.Matrix) *mat.Dense {
	ret := &mat.Dense{}
	ret.Mul(a, b)
	return ret
}

func _add(a, b mat.Matrix) *mat.Dense {
	ret := &mat.Dense{}
	ret.Add(a, b)
	return ret
}

func _sub(a, b mat.Matrix) *mat.Dense {
	ret := &mat.Dense{}
	ret.Sub(a, b)
	return ret
}

package track

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
type SimplifyOption struct {
	Degree float64 `protobuf:"fixed64,1,opt,name=degree,proto3" json:"degree,omitempty"`
}

type Simplify struct {
}

func NewSimplify() *Simplify {
	return &Simplify{}
}

func (d *Simplify) Exec(ops *SimplifyOption, coords ...TrackPointer) []TrackPointer {
	return DouglasPeucker(d._transEpsilon(int(ops.Degree)), coords...)
}

func (d *Simplify) _transEpsilon(level int) float64 {
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

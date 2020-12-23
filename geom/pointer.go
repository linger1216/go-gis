package geom

type Pointer interface {
	ID() string
	Point() *LngLat
	Timestamp() int64
}

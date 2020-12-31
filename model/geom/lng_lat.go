package geom

type LngLat struct {
	Longitude float64
	Latitude  float64
	Altitude  float64
}

func (l *LngLat) ID() string {
	return ""
}

func (l *LngLat) Point() *LngLat {
	return l
}

func (l *LngLat) Timestamp() int64 {
	return 0
}

func NewLngLat(coordinates ...float64) *LngLat {
	length := len(coordinates)
	switch length {
	case 0, 1:
		return nil
	case 2:
		return &LngLat{Longitude: coordinates[0], Latitude: coordinates[1]}
	case 3:
		return &LngLat{
			Longitude: coordinates[0],
			Latitude:  coordinates[1],
			Altitude:  coordinates[2],
		}
	default:
		return nil
	}
}

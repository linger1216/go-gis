package geom

import "math"

const (
	// According to Wikipedia, the Earth's radius is about 6,371,004 m
	EarthRadius = 6371000
)

func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

func Distance(lng1, lat1, lng2, lat2 float64) float64 {
	// convert to radians
	// must cast radius as float to multiply later
	var la1, ln1, la2, ln2, r float64
	la1 = lat1 * math.Pi / 180
	ln1 = lng1 * math.Pi / 180
	la2 = lat2 * math.Pi / 180
	ln2 = lng2 * math.Pi / 180

	r = EarthRadius // Earth radius in METERS

	// calculate
	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(ln2-ln1)

	return 2 * r * math.Asin(math.Sqrt(h))
}

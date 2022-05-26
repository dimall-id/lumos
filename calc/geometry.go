package calc

import "math"

func DistanceBetween (lat1 float64, lon1 float64, lat2 float64, lon2 float64) float64 {
	R := 6371.0
	dLat := deg2Rad(lat2 - lat2)
	dLon := deg2Rad(lon2 - lon1)
	a := (math.Sin(dLat/2) * math.Sin(dLat/2)) + (math.Cos(deg2Rad(lat1)) * math.Cos(deg2Rad(lat2)) * math.Sin(dLon/2) * math.Sin(dLon/2))
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1 - a))
	d := R * c
	return d
}

func deg2Rad (deg float64) float64 {
	return deg * (math.Pi/180.0)
}
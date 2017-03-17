package geo

// Maybe move to https://github.com/paulmach/go.geo ?

// Formulas from http://www.movable-type.co.uk/scripts/latlong.html
// Quadtree walking: https://www.cs.umd.edu/class/spring2008/cmsc420/L17-18.QuadTrees.pdf

import "math"

const (
	kmtomiles = float64(0.621371192)
	earthRadiusKM = float64(6371)
	kKMPerNauticalMile = float64(1.852)

	KFeetPerKM = float64(3280.8399)
	KNauticalMilePerKM = float64(0.539957)
)

func NM2KM (nm float64) float64 { return nm / KNauticalMilePerKM }
// NM per hour to meters per second
func NMph2mps(f float64) float64 { return f * 0.514444444 }

// The haversine formula will calculate the spherical distance as the crow flies 
// between lat and lon for two given points, in km
func haversine(lonFrom float64, latFrom float64, lonTo float64, latTo float64) float64 {
	var deltaLat = (latTo - latFrom) * (math.Pi / 180)
	var deltaLon = (lonTo - lonFrom) * (math.Pi / 180)
	var a = math.Sin(deltaLat / 2) * math.Sin(deltaLat / 2) + 
		math.Cos(latFrom * (math.Pi / 180)) * math.Cos(latTo * (math.Pi / 180)) *
		math.Sin(deltaLon / 2) * math.Sin(deltaLon / 2)
	var c = 2 * math.Atan2(math.Sqrt(a),math.Sqrt(1-a))
	return earthRadiusKM * c
}

// Computes the intitial bearing to head from 1 to 2 in a 'straight line great circle'
func forwardAzimuth(lon1,lat1 float64, lon2,lat2 float64) float64 {
	lon1R,lat1R := (lon1 * (math.Pi / 180.0)), (lat1 * (math.Pi / 180.0)) // To radians
	lon2R,lat2R := (lon2 * (math.Pi / 180.0)), (lat2 * (math.Pi / 180.0))

	y := math.Sin(lon2R-lon1R) * math.Cos(lat2R)
	x := ( math.Cos(lat1R) * math.Sin(lat2R) ) -
 		   ( math.Sin(lat1R) * math.Cos(lat2R) * math.Cos(lon2R-lon1R) )

	bearing := math.Atan2(y, x) * (180.0 / math.Pi)
	return math.Mod(bearing+360.0, 360.0)
}

// Computes a new position, given a start position, a heading, and a distance. Expects
// latlongs in degrees.
func move(lon1, lat1, bearing, distanceKM float64) (float64,float64) {
	lon1R := lon1 * (math.Pi / 180.0)
	lat1R := lat1 * (math.Pi / 180.0) // To radians

	bearing = math.Mod(bearing+360, 360) * (math.Pi / 180.0)
	
	lat2R := math.Asin(
		math.Sin(lat1R) * math.Cos(distanceKM/earthRadiusKM) +
	  math.Cos(lat1R) * math.Sin(distanceKM/earthRadiusKM) * math.Cos(bearing) )
	
	lon2R := lon1R + math.Atan2(
		math.Sin(bearing) * math.Sin(distanceKM/earthRadiusKM) * math.Cos(lat1R),
		math.Cos(distanceKM/earthRadiusKM) - math.Sin(lat1R) * math.Sin(lat2R) )

	lon2 := lon2R * (180.0 / math.Pi)
	lat2 := lat2R * (180.0 / math.Pi)
	
	lon2 = math.Mod((lon2+540.0), 360.0) - 180.0 // normalize to [-180,180]

	return lon2,lat2
}


// {{{ -------------------------={ E N D }=----------------------------------

// Local variables:
// folded-file: t
// end:

// }}}

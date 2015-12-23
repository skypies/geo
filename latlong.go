package geo

import(
	"fmt"
	"math"
)

type Latlong struct {
	Lat, Long float64
}
func (ll Latlong)String() string { return fmt.Sprintf("(%.4f,%.4f)", ll.Lat, ll.Long) }

// ApproxDist treats the latlongs as y,x coords, and returns the 'latlong dist'
func (from Latlong)LatlongDist(to Latlong) float64 {
	x,y := (to.Long - from.Long), (to.Lat - from.Lat)
	return math.Sqrt(x*x + y*y)
}

// Dist is the great-circle distance, in KM
func (from Latlong)Dist(to Latlong) float64 {
	return haversine(from.Long,from.Lat,  to.Long,to.Lat)
}
func (from Latlong)DistNM(to Latlong) float64 {
	return from.Dist(to) * kNauticalMilePerKM
}

func (from Latlong)BearingTowards(to Latlong) float64 {
	return forwardAzimuth(from.Long,from.Lat,  to.Long,to.Lat)
}

func (at Latlong)MapsUrl() string {
	return fmt.Sprintf("https://www.google.com/maps/@%.6f,%.6f,9z", at.Lat, at.Long)
}

// We have the two right angle sides of the triangle !
func (from Latlong)Dist3(to Latlong, altitude float64) float64 {
	horizDist := haversine(from.Long,from.Lat,  to.Long,to.Lat)
	vertDist  := altitude / 3280.84
	return math.Sqrt(horizDist*horizDist + vertDist*vertDist)
}

func (from Latlong)InterpolateTo(to Latlong, ratio float64) Latlong {
	interpFunc := func(from,to float64) float64 { return from + (to-from)*ratio }

	return Latlong{
		Lat: interpFunc(from.Lat, to.Lat),
		Long: interpFunc(from.Long, to.Long),
	}
}

// {{{ -------------------------={ E N D }=----------------------------------

// Local variables:
// folded-file: t
// end:

// }}}

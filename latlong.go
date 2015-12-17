package geo

import(
	"fmt"
	"math"
)

// {{{ type Latlong

type Latlong struct {
	Lat, Long float64
}
func (ll Latlong)String() string { return fmt.Sprintf("(%.4f,%.4f)", ll.Lat, ll.Long) }

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

// }}}
// {{{ type LatlongBox

const (
	kKmPerLatitudeDegreeAtSFO  = 111.2   // (36,-122)->(37,-122) == 111 KM (heading north)
	kKmPerLongitudeDegreeAtSFO = 88.08   // (36,-122)->(36,-121) ==  88 KM (heading east)
)

type LatlongBox struct {
	SW, NE Latlong
}
func (box LatlongBox)String() string { return fmt.Sprintf("%s-%s", box.SW, box.NE) }

// Returns a box, centred on ll, that is of size (width,height)
func (ll Latlong)Box(widthKm,heightKm float64) LatlongBox {
	// This is a hack; these constants are only valid quite close to SFO itself.
	latOffset := (heightKm / kKmPerLatitudeDegreeAtSFO) / 2.0
	longOffset := (widthKm / kKmPerLongitudeDegreeAtSFO) / 2.0
	return LatlongBox{
		SW: Latlong{ll.Lat-latOffset, ll.Long-longOffset},
		NE: Latlong{ll.Lat+latOffset, ll.Long+longOffset},
	}
}

func (from Latlong)BoxTo(to Latlong) LatlongBox {
	box := LatlongBox{}
	if (from.Lat < to.Lat) {
		box.SW.Lat, box.NE.Lat = from.Lat, to.Lat
	} else {
		box.SW.Lat, box.NE.Lat = to.Lat, from.Lat
	}
	if (from.Long < to.Long) {
		box.SW.Long, box.NE.Long = from.Long, to.Long
	} else {
		box.SW.Long, box.NE.Long = to.Long, from.Long
	}
	return box
}

func (box LatlongBox)Contains(l Latlong) bool {
	if (l.Long < box.SW.Long) { return false }
	if (l.Lat  < box.SW.Lat ) { return false }
	if (l.Long > box.NE.Long) { return false }
	if (l.Lat  > box.NE.Lat ) { return false }
	return true
}

// }}}

// {{{ -------------------------={ E N D }=----------------------------------

// Local variables:
// folded-file: t
// end:

// }}}

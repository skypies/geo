package geo

import(
	"fmt"
)

const (
	kKmPerLatitudeDegreeAtSFO  = 111.2   // (36,-122)->(37,-122) == 111 KM (heading north)
	kKmPerLongitudeDegreeAtSFO = 88.08   // (36,-122)->(36,-121) ==  88 KM (heading east)
)

type LatlongBox struct {
	SW, NE Latlong
}
func (box LatlongBox)String() string { return fmt.Sprintf("%s-%s", box.SW, box.NE) }

// Derive the other two corners on demand
func (box LatlongBox)SE() Latlong { return Latlong{Lat:box.SW.Lat , Long:box.NE.Long} }
func (box LatlongBox)NW() Latlong { return Latlong{Lat:box.NE.Lat , Long:box.SW.Long} }

func (box LatlongBox)LongWidth() float64 { return box.NE.Long - box.SW.Long }
func (box LatlongBox)LatHeight() float64 { return box.NE.Lat - box.SW.Lat }

func (box LatlongBox)Center() Latlong {
	return Latlong{
		Lat: (box.SW.Lat + box.NE.Lat) / 2.0,
		Long: (box.SW.Long + box.NE.Long) / 2.0,
	}
}

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

func (box LatlongBox)Contains(pos Latlong) bool {
	if (pos.Long < box.SW.Long) { return false }
	if (pos.Lat  < box.SW.Lat ) { return false }
	if (pos.Long > box.NE.Long) { return false }
	if (pos.Lat  > box.NE.Lat ) { return false }
	return true
}

// Enclose increases the sixe of the box to include the point, if it doesn't fit
func (box *LatlongBox)Enclose(pos Latlong) {
	if (pos.Long < box.SW.Long) { box.SW.Long = pos.Long }
	if (pos.Lat  < box.SW.Lat ) { box.SW.Lat = pos.Lat }
	if (pos.Long > box.NE.Long) { box.NE.Long = pos.Long }
	if (pos.Lat  > box.NE.Lat ) { box.NE.Lat = pos.Lat }
}

func (box LatlongBox)LatRange() Float64Range { return Float64Range{box.SW.Lat, box.NE.Lat} }
func (box LatlongBox)LongRange() Float64Range { return Float64Range{box.SW.Long, box.NE.Long} }

// We take some liberties with the outcome value; don't rely on it except for .IsDisjoint()
// float is the fraction of b1 that overlaps with b2
func (b1 LatlongBox)OverlapsWith(b2 LatlongBox) (OverlapOutcome,float64) {
	latDisp := RangeOverlap(b1.LatRange(), b2.LatRange())
	longDisp := RangeOverlap(b1.LongRange(), b2.LongRange())
	
	if latDisp.IsDisjoint() || longDisp.IsDisjoint() {
		return DisjointR2ComesAfter, 0.0
	}

	// Really, we should figure out the actual amount of overlap; but for now, hardwire 1.0
	return OverlapR2StraddlesStart, 1.0
}

// {{{ -------------------------={ E N D }=----------------------------------

// Local variables:
// folded-file: t
// end:

// }}}

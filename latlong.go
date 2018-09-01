package geo

import(
	"fmt"
	"math"
	"regexp"
	"strconv"
)

type Latlong struct {
	Lat  float64
	Long float64
}
func (ll Latlong)String() string { return fmt.Sprintf("(%.4f,%.4f)", ll.Lat, ll.Long) }

// Recognizes some common formats:
//   [36.7415306, -121.8942333] - google maps style full decimals
//   [36°57'02.96"N, 121°57'09.62"W] - traditional degrees / minutes / seconds
//   [365702.96N / 1215709.62W] - FAA style concatenated ("DEG"+"MIN"+"SEC.00")
//   [37-10-35.680N / 122-00-29.950W] - from fltplan.com
// returns a nil latlong upon parse failure
func NewLatlong(in string) Latlong {
	coordRe := `[-0-9.°'\"NESW]{5,16}`
	re := regexp.MustCompile(`^\s*(`+coordRe+`)[,/\s]+(`+coordRe+`)\s*$`)

	if match := re.FindStringSubmatch(in); match != nil && len(match)==3 {
		return Latlong{
			Lat: parseCoord(match[1]),
			Long:parseCoord(match[2]),
		}
	}

	return Latlong{}
}

func parseCoord(in string) float64 {
	decimalRe   := regexp.MustCompile(`^-?\d{1,3}\.\d{3,9}$`)
	concatRe    := regexp.MustCompile(`^(\d{1,3})(\d\d)(\d\d(?:\.\d+))([NEWS])$`)
	degMinSecRe := regexp.MustCompile(`^(\d{1,3})[°'\"\s]+(\d{2})[°'\"\s]+(\d{2}(?:\.\d+)?)[°'\"\s]+([NEWS])$`)

	// Given degrees, minutes, seconds, and a compass dir; return a decimal coord
	dms2dec := func(strs []string) float64 {
		d,_ := strconv.ParseInt(strs[0], 10, 64)
		m,_ := strconv.ParseInt(strs[1], 10, 64)
		s,_ := strconv.ParseFloat(strs[2], 64)
		dir := float64(1.0)
		if strs[3] == "W" || strs[3] == "S" { dir *= -1 }
		return (float64(d) + float64(m)/60.0 + s/3600.0) * dir
	}

	if decimalRe.MatchString(in) {
		coord,_ := strconv.ParseFloat(in, 64) // no errors ;)
		return coord
	} else if match := degMinSecRe.FindStringSubmatch(in); match != nil {//&& len(match) == 5 {
		return dms2dec(match[1:])
	} else if match := concatRe.FindStringSubmatch(in); match != nil && len(match) == 5 {
		return dms2dec(match[1:])
	}
	return 0
}

type LatlongSlice []Latlong

// We often treat latlong as a simple (x,y) space. We take Long as x, to make horiz/vert look normal
func (ll Latlong)x() float64 { return ll.Long }
func (ll Latlong)y() float64 { return ll.Lat }

// The square of the distance; useful for deciding closest approach
func (from Latlong)LatlongDistSq(to Latlong) float64 {
	x,y := (to.Long - from.Long), (to.Lat - from.Lat)
	return x*x + y*y
}

// ApproxDist treats the latlongs as y,x coords, and returns the 'latlong dist'
func (from Latlong)LatlongDist(to Latlong) float64 {
	return math.Sqrt(from.LatlongDistSq(to))
}

// Maybe it's safe to compare a nil-float64 directly to 0.0, and this is just a paranoid bodge.
func (ll Latlong)IsNil() bool {
	return math.Abs(ll.Lat)<0.01 && math.Abs(ll.Long)<0.01
}

// This probably isn't what you want
func (from Latlong)ExactlyEqual(to Latlong) bool {
	return from.Lat == to.Lat && from.Long == to.Long
}

func (from Latlong)Equal(to Latlong) bool {
	return floatEquals(from.Lat, to.Lat) && floatEquals(from.Long, to.Long)
}

const EPSILON float64 = 0.0000001
func floatEquals(a, b float64) bool {
	if ((a - b) < EPSILON && (b - a) < EPSILON) {
		return true
	}
	return false
}

// Dist is the great-circle distance, in KM
func (from Latlong)Dist(to Latlong) float64 { return from.DistKM(to) }
func (from Latlong)DistKM(to Latlong) float64 {
	return haversine(from.Long,from.Lat,  to.Long,to.Lat)
}
func (from Latlong)DistNM(to Latlong) float64 {
	return from.DistKM(to) * KNauticalMilePerKM
}

func (from Latlong)BearingTowards(to Latlong) float64 {
	return forwardAzimuth(from.Long,from.Lat,  to.Long,to.Lat)
}

func (from Latlong)MoveKM(heading, distanceKM float64) Latlong {
	long,lat := move(from.Long, from.Lat, heading, distanceKM)
	return Latlong{Lat:lat, Long:long}
}
func (from Latlong)MoveNM(heading, distanceNM float64) Latlong {
	return from.MoveKM(heading, distanceNM * KNauticalMilePerKM)
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

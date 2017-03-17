package geo
// go test -v github.com/skypies/geo

import (
	"fmt"
	"math"
	"testing"
)

const (
	kLatSFO  = 37.6188172
	kLongSFO = -122.3754281
	kLatSJC  = 37.3639472
	kLongSJC = -121.9289375
)

var (
	// All these data have (long,lat) pairs - don't get confused !
	bearings = [][]float64{
		{-120, 36,         -120, 37,            00.00},  // Due north
		{-120,  0,         -110,  0,            90.00},  // Due east (at equator)
		{-120, 37,         -120, 36,           180.00},  // Due south
		{-120,  0,         -130,  0,           270.00},  // Due west (at equator)
		{kLongSFO, kLatSFO, kLongSJC,kLatSJC,  125.596}, // SFO to SJC
	}

	distances = [][]float64{
		{kLongSJC,kLatSJC, kLongSFO, kLatSFO,  48.528},  // 48 km
	}

	toSFO = [][]float64{
		// long,lat,        distNM,    bearing
		{kLongSJC,kLatSJC,  26.202843, 305.868158},
	}

	SFOClassB = [][]float64{
		// long, lat,             inRange, floor, ceiling
		{-122.291730,  38.306015, -1.0 }, // Napa, way north, not in range
		{-121.994670,  37.203240, -1.0 }, // Lexington reservoir; just outside
		{-122.020241,  36.986304, -1.0 }, // Santa Cruz, not in range
		{-122.022386,  37.273368,  1.0, 8000, 10000 },  // 100/80, Saratoga
		{-122.153278,  37.278354,  1.0, 6000, 10000 },  // 100/60, Long Ridge open space preserve
		{-122.204562,  37.418809,  1.0, 4000, 10000 },  // 100/40, SLAC accelerator
		{-122.438751,  37.479148,  1.0, 4000, 10000 },  // 100/40, Half Moon bay
	}

	// These *better* tests order values as lat,long  !!
	latlongBoxes = [][]float64{
		// lat,   long,   width, height, distance we expect (all distances in KM)
		{kLatSFO, kLongSFO, 50.0,  0.0,  49.998},
		{kLatSFO, kLongSFO,  0.0, 50.0,  49.998},
		{kLatSFO, kLongSFO, 50.0, 50.0,  70.707},
	}

	containsBox = LatlongBox{ SW: Latlong{36,-122}, NE: Latlong{40,-118} }
	containsPoints = [][]float64{
		// lat, long, [>0 if true]
		{ 36, -122,    1.0},
		{ 37, -121,    1.0},
		{ 40, -118,    1.0},
		{ 35, -120,   -1.0},
		{ 41, -120,   -1.0},
		{ 38, -123,   -1.0},
		{ 38, -117,   -1.0},
	}

	closestDistance = [][]float64{
		// line coords,       point,      dist (in KM; ~90KM per unit)
		{ 36,-120, 37,-120,   36.5,-120,    0.00},
		{ 37,-120, 36,-120,   36.5,-121,   89.3844},
		{ 36,-120, 37,-119,   36.5,-119.4,  7.1325},
	}

	distalongline = [][]float64{
		// long,lat,  long,lat,   long,lat,   dist
		{  -120,35,   -118,35,    -121,35,   -0.50 }, // Horizontal line
		{  -120,35,   -118,35,    -120,35,    0.00 },
		{  -120,35,   -118,35,    -118.5,35,  0.75 },
		{  -120,35,   -118,35,    -116,35,    2.00 },
		{  -120,35,   -118,35,    -119,37,    0.50 }, // Point doesn't lie on line - lies above it

		{  -120,35,   -120,39,    -120,35,    0.00 }, // Vertical line
		{  -120,35,   -120,39,    -120,38,    0.75 },
		{  -120,35,   -120,39,    -120,47,    3.00 },

		{  -120,35,   -118,37,    -121,34,   -0.50 }, // at 45 degrees
		{  -120,35,   -118,37,    -120,35,    0.00 },
		{  -120,35,   -118,37,    -119,36,    0.50 },
		{  -120,35,   -118,37,    -117,38,    1.50 },

	}
)

func floatEq(x,y float64) bool { return (math.Abs(x - y) <= 0.001) }

// http://www.movable-type.co.uk/scripts/latlong.html
func TestDist(t *testing.T) {
	for i,vals := range distances {
		from,to := Latlong{vals[1], vals[0]}, Latlong{vals[3], vals[2]}
		actual := from.Dist(to)
		if (math.Abs(actual - vals[4]) > 0.001) {
			t.Errorf("[%d] Haversine said '%f', given %v", i, actual, vals)
		}
	}
}

func TestLatlongbox(t *testing.T) {
	for i,vals := range latlongBoxes {
		centre := Latlong{vals[0], vals[1]}
		box := centre.Box(vals[2], vals[3])
		actual := box.SW.Dist(box.NE)
		if !floatEq(actual, vals[4]) {
			t.Errorf("[%d] Latlong corner to corner was '%f', given %v", i, actual, vals)
		}
	}
}

func TestForwardAzimuth(t *testing.T) {
	for i,vals := range bearings {
		from,to := Latlong{vals[1], vals[0]}, Latlong{vals[3], vals[2]}
		actual := from.BearingTowards(to)
		if (math.Abs(actual - vals[4]) > 0.001) {
			t.Errorf("[%d] Azimuth said '%f', given %v", i,actual, vals)
		}
	}
}

func TestLatlongBoxContains(t *testing.T) {
	for i,vals := range containsPoints {
		p := Latlong{vals[0], vals[1]}
		actual := containsBox.Contains(p)
		if actual != (vals[2]>0) {
			t.Errorf("[%d] Contains gave %v, expected %v", i, actual, (vals[2]>0))
		}
	}
}
/*
func TestToSFO(t *testing.T) {
	for i,vals := range toSFO {
		from := Latlong{vals[1], vals[0]}
		actualDist,actualBearing := from.ToSFO()
		if (math.Abs(actualDist - vals[2]) > 0.001) {
			t.Errorf("[%d] ToSFO gave distNM of '%f', wanted %f", i,actualDist, vals[2])
		}
		if (math.Abs(actualBearing - vals[3]) > 0.001) {
			t.Errorf("[%d] ToSFO gave bearing of '%f', wanted %f", i, actualBearing, vals[3])
		}
	}	
}

func TestSFOClassBRange(t *testing.T) {
	for i,vals := range SFOClassB {
		floor, ceil, inRange := SFOClassBRange(Latlong{vals[1], vals[0]})
		if (inRange != (vals[2] > 0.0)) {
			t.Errorf("[%d] inRange was %v, given %v", i, inRange, vals)
		} else if (inRange) {
			if (math.Abs(floor - vals[3]) > 0.001) {
				t.Errorf("[%d] floor was %f, given %v", i, floor, vals)
			}
			if (math.Abs(ceil - vals[4]) > 0.001) {
				t.Errorf("[%d] ceil was %f, given %v", i, ceil, vals)
			}
		}
	}
}
*/
func TestClosestDistance(t *testing.T) {
	for i,vals := range closestDistance {
		lFrom,lTo,pos := Latlong{vals[0],vals[1]}, Latlong{vals[2],vals[3]}, Latlong{vals[4],vals[5]}
		line := lFrom.BuildLine(lTo)
		//perp := line.PerpendicularTo(pos)
		//closestPoint,parallel := line.intersectByLineEquations(perp)
		dist := line.ClosestDistance(pos)
		if math.Abs(dist - vals[6]) > 0.001 {
			t.Errorf("[%d] dist was %f, expected %f", i, dist, vals[6])
		}
		//fmt.Printf("Line:%s, pos    :%s\nPerp:%s, closest:%s (parallel=%v)\ndist=%.2f\n--\n",
		//	line, pos, perp, closestPoint, parallel, dist)
	}
}

func TestDistAlongLine(t *testing.T) {
	for i,vals := range distalongline {
		lFrom,lTo,pos := Latlong{vals[0],vals[1]}, Latlong{vals[2],vals[3]}, Latlong{vals[4],vals[5]}
		line := lFrom.BuildLine(lTo)
		actual := line.DistAlongLine(pos)
		expected := vals[6]
		if math.Abs(actual-expected) > 0.001 {
			t.Errorf("[%d] distalongline was %f, expected %f", i, actual, expected)
		}
		// fmt.Printf("Line:%s, pos:%s, dist:%.3f\n", line, pos, actual)
	}
}



func degstr(f float64) string {
	deg := int(f)
	minf := math.Abs(float64(deg) - f) * 60.0
	min := int(minf)
	sec := (minf - float64(min)) * 60.0

	return fmt.Sprintf("% 3d°%02d'%02.0f\"", deg, min, sec)
}
func (p Latlong)DegStr() string {
	return fmt.Sprintf("(%sN, %sE)", degstr(p.Lat), degstr(p.Long))
}

func (p Latlong)Eq(q Latlong) bool {
	if math.Abs(p.Lat - q.Lat) > 0.000001 { return false }
	if math.Abs(p.Long - q.Long) > 0.000001 { return false }
	return true
}

//	func move(lon1, lat1, heading, distanceKM float64) (lon2, lat2 float64) {
func TestMove(t *testing.T) {
	tests := []struct{
		From Latlong
		Heading, Distance float64
		Expected Latlong
	}{
		// Simple tests - north, east, south, then west.
		// Sanity checked vs. http://www.movable-type.co.uk/scripts/latlong.html
		{Latlong{36,-120},    0, 100, Latlong{36.899322, -120}},
		{Latlong{36,-120},   90, 100, Latlong{35.994872, -118.888426}},
		{Latlong{36,-120},  180, 100, Latlong{35.100678, -120}},
		{Latlong{36,-120},  270, 100, Latlong{35.994872, -121.111574}},
	}

	for i,test := range tests {
		actual := test.From.MoveKM(test.Heading, test.Distance)
		if !actual.Eq(test.Expected) {
			t.Errorf("test [% 2d], %s+(%3.0fdeg,%.0fKM) - expected %s, got %s\n",
				i, test.From, test.Heading, test.Distance, test.Expected, actual)
		}
	}
}


// {{{ -------------------------={ E N D }=----------------------------------

// Local variables:
// folded-file: t
// end:

// }}}

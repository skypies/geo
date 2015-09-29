package geo

import "fmt"

type Cylinder struct {
	EndDistanceNM int  // Nautical Miles. Start distance is end of prev cylinder (or origin)
	Floor      int     // In hundreds of feet
	Ceil       int     // In hundreds of feet
}

// Each sector is a pie wedge, with consistent floor/ceil cylinders
type ClassBSector struct {
	StartBearing   int  // Magnetic bearing (-13.68 to get to magnetic @ SFO)
	EndBearing     int
	Steps        []Cylinder // Ordered by asc DistanceNM
}

type ClassBMap struct {
	Sectors  []ClassBSector  // Must be ordered by asc StartBearing, and support wraparound
	Center   Latlong
	Name     string
}

func (m ClassBMap) Walk(distNM, bearing float64) (floor,ceil int, inRange bool) {
	inRange = false
	// Walk the sectors until we find the first one which contains our bearing
	// TODO: funky wraparound thing
	for _,sector := range m.Sectors {
		if bearing < float64(sector.EndBearing) {
			// Now, walk the Cylinders and find first one that contains our distance
			for _,cyl := range sector.Steps {
				if distNM < float64(cyl.EndDistanceNM) {
					floor, ceil, inRange = cyl.Floor, cyl.Ceil, true
					return // We found our class B limits !
				}
			}
		}
		return // We are past the outer limits of this sector's cylinders
	}

	panic(fmt.Sprintf("Bad ClassBMap, we fell off the end, given bearing=%f", bearing))
	return
}

func (m ClassBMap)ClassBRange(pos Latlong) (floor,ceil float64, inRange bool) {
	floor,ceil,inRange = 0.0, 0.0, false
	distNM := pos.DistNM(m.Center)
	bearing := pos.BearingTowards(m.Center)

	var f,c int
	f,c,inRange = m.Walk(distNM, bearing)
	if inRange {
		floor = float64(f) * 100.0
		ceil = float64(c) * 100.0
	}
	return
}

// A trackpoint, annotated with ClassB analysis
type TPClassBAnalysis struct {
	// The verdict
	WithinRange         bool    // If we're not within range, the rest has no meaning.
	VerticalDisposition int     // <0 below; =0 within; >0 above
	BelowBy             float64 // If below, by how many feet
	Reasoning           string  // Explanation of stuff

	// Handy data to have around later
	Floor,Ceil          float64 // The Class B space the point was in (0 if not in space)
	DistNM              float64 // Seeing as we've calculated it :)

	AllowThisPoint      bool    // If true, this point is not a violation, regardless of data
	
	// Old value, for posterity
	OldOutcome          int
}
func (a TPClassBAnalysis)IsViolation() bool {
	if a.AllowThisPoint { return false }
	if !a.WithinRange { return false }
	return a.VerticalDisposition < 0
}

// Note: we're ignoring speed limits, and violation detection on a per-point basis; violations
// are detected by a state machine that walks successive outcomes from this function.
func (m ClassBMap)ClassBPointAnalysis(pos Latlong, speed float64, altitude float64, o *TPClassBAnalysis) {
	distNM := pos.DistNM(m.Center)
	bearing := pos.BearingTowards(m.Center)
	o.DistNM = distNM

	o.Reasoning = fmt.Sprintf("** ClassB analysis: aircraft at %s, %.0f kt, %.0f feet\n",
		pos,speed,altitude)
	o.Reasoning += fmt.Sprintf("* Distance to %s in NM: %.1f; bearing towards %s: %.1f\n",
		m.Name, distNM, m.Name, bearing)

	o.Floor,o.Ceil,o.WithinRange = m.ClassBRange(pos)
	
	if !o.WithinRange {
		o.Reasoning += "* not in range; too far away from "+m.Name+"\n"
		return
	}

	limitStr := fmt.Sprintf("%d/%d", int(o.Ceil/100.0), int(o.Floor/100.0))
	o.Reasoning += fmt.Sprintf("* In <b>%s</b> space, at <b>%.0f</b> feet\n", limitStr, altitude)
	
	if (altitude > o.Ceil) {
		o.VerticalDisposition = 1
		o.Reasoning += "* above class B ceiling\n"
		
	} else if (altitude > o.Floor-201) {  // Allow 200' of wriggle room
		o.VerticalDisposition = 0
		o.Reasoning += "* within class B height range\n"
		
	} else {
		o.VerticalDisposition = -1
		o.Reasoning += "* below class B floor\n"
		o.BelowBy = o.Floor - altitude
	}

	return
}

// {{{ -------------------------={ E N D }=----------------------------------

// Local variables:
// folded-file: t
// end:

// }}}

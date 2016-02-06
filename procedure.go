package geo

import "fmt"

type Waypoint struct {
	FixName     string
	Latlong     // embed - this needs to be populated via a fix table !
	MinAltitude float64
	MaxAltitude float64
	MaxAirspeed float64

	ConsistentlyFlown bool // Can we assume that aircraft on this procedure will fly through here ?
}

// A Procedure is a set of waypoints, with altitude / speed guidance, for arrival/departure
type Procedure struct {
	Name        string
	Departure   bool     // If false, is arrival
	Airport     string   // Where we are arriving or departing from
	Waypoints []Waypoint
}

func (p Procedure)String() string {
	str := fmt.Sprintf("%s:-\n", p.Name)
	for _,wp := range p.Waypoints {
		str += fmt.Sprintf("  %s (%5d-%5d) %3dk", wp.FixName, wp.MinAltitude, wp.MaxAltitude,
			wp.MaxAirspeed)
		if wp.ConsistentlyFlown { str += " [**]" }
		str += "\n"
	}
	return str
}

func (p *Procedure)Populate(fixes map[string]Latlong) {
	for i,_ := range p.Waypoints {
		p.Waypoints[i].Latlong = fixes[p.Waypoints[i].FixName]
	}
}

// ComparisonLines returns the segments of the procedure that a track's boxes need to intersect and
// align with.
func (p Procedure)ComparisonLines() []LatlongLine {
	ret := []LatlongLine{}
	for i,_ := range p.Waypoints[1:] {
		if p.Waypoints[i].ConsistentlyFlown && p.Waypoints[i+1].ConsistentlyFlown {
			ret = append(ret, p.Waypoints[i].BuildLine(p.Waypoints[i+1].Latlong))
		}
	}
	return ret
}

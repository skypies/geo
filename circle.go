package geo

import "fmt"

type LatlongCircle struct {
	Latlong
	RadiusKM float64 // In kilometers
}

func (c LatlongCircle)String() string {
	return fmt.Sprintf("%s, rad=%.1f", c.Latlong.String(), c.RadiusKM)
}

func (pos Latlong)Circle(radius float64) LatlongCircle {
	return LatlongCircle{pos, radius}
}

// Implement GeoRestricter interface
func (c LatlongCircle)LookForExit() bool { return true }
func (c LatlongCircle)IntersectsLine(l LatlongLine) bool {
	// http://math.stackexchange.com/questions/228841/how-do-i-calculate-the-intersections-of-a-straight-line-and-a-circle
	return false // Implement me
}
func (c LatlongCircle)IntersectsAltitude(alt int64) bool {
	return false // Implement me, too
}


func (c LatlongCircle)IntersectsLineDeb(l LatlongLine) (bool,string) {return c.IntersectsLine(l),""}

// Implement Region interface (defunct ?)
func (c LatlongCircle)ContainsPoint(pos Latlong) bool {
	return (c.DistKM(pos) < c.RadiusKM)
}
func (c LatlongCircle)IntersectsBox(b2 LatlongBox) bool {
	return false // Implement me
}

// Implement MapRenderer interface
func (c LatlongCircle)ToLines() []LatlongLine { return []LatlongLine{} }
func (c LatlongCircle)ToCircles() []LatlongCircle { return []LatlongCircle{c} }

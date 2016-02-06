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

func (c LatlongCircle)ContainsPoint(pos Latlong) bool {
	return (c.DistKM(pos) < c.RadiusKM)
}
func (c LatlongCircle)IntersectsLine(l LatlongLine) bool {
	return false // Implement me
}
func (c LatlongCircle)IntersectsBox(b2 LatlongBox) bool {
	return false // Implement me
}
func (c LatlongCircle)ToLines() []LatlongLine { return []LatlongLine{} }
func (c LatlongCircle)ToCircles() []LatlongCircle { return []LatlongCircle{c} }

package geo

import(
	"fmt"
	"math"
)

const kLineSnapKM = 0.3  // How far a trackpoint can be from a line, and still be on that line
const KLineSnapKM = 0.3  // How far a trackpoint can be from a line, and still be on that line

// Always using two anchor points, and then derive the equation of the line: y = m.x + b
type LatlongLine struct {
	From,To  Latlong
	m,b      float64

	I,J      int  // index values for the two points used to build this line
}
func (line LatlongLine)String() string {
	return fmt.Sprintf("[y=%.2f.x + %.2f] (%.3f,%.3f)->(%.3f,%.3f) [i=%d,j=%d]",
		line.m, line.b, line.From.x(), line.From.y(), line.To.x(), line.To.y(), line.I, line.J)
}

// {{{ calcM, calcB

// Helpers for line construction
func calcM(p1,p2 Latlong) float64 { return (p2.y() - p1.y()) / (p2.x() - p1.x()) }
func calcB(m float64, p Latlong) float64 {	
	// Given a gradient(m) and a point, work out b (the value of y when x==0)
	
	if math.IsInf(m,0) { return math.NaN() } // Equation of line does not apply for vertical lines

	// y=m.x+b; so b=y-m.x for both points (x,y)
	return p.y() - (m * p.x())
}

// }}}

// {{{ latlong.LineTo (BuildLine)

func (from Latlong)LineTo(to Latlong) LatlongLine {
	m := calcM(from,to)
	return LatlongLine{
		From: from,
		To: to,
		m: m,
		b: calcB(m,to),
	}
}

func (from Latlong)BuildLine(to Latlong) LatlongLine { return from.LineTo(to) }

// }}}

// {{{ l.x, l.y, l.Box

// Apply equation of line: y=mx+b
func (line LatlongLine)y(x float64) float64 { return line.m * x + line.b }
func (line LatlongLine)x(y float64) float64 { return (y - line.b) / line.m }

func (l LatlongLine)Box() LatlongBox { return l.From.BoxTo(l.To) }

func (l LatlongLine)IsVertical()   bool { return math.IsInf(l.m,0) }
func (l LatlongLine)IsDegenerate() bool {
	return l.From.Lat==l.To.Lat && l.From.Long==l.To.Long
}

// }}}
// {{{ l.intersectByLineEqutions

// This function uses the m,b line constants.
// If either line is vertical, we use l.From anchor point; there is no need for a l.To point.
// The returned bool is true if lines were parallel.
// https://en.wikipedia.org/wiki/Line%E2%80%93line_intersection#Given_the_equations_of_the_lines
func (l1 LatlongLine)intersectByLineEquations(l2 LatlongLine) (Latlong, bool) {
	if l1.m == l2.m { return Latlong{}, true } // same slope; are parallel

	// 1. y=ax+c   2. y=bx+d
	a,c := l1.m,l1.b
	b,d := l2.m,l2.b

	if math.IsInf(a,0) {
		// l1 is vertical; x is fixed, take it from the anchor point. Find the point on l2 for that x.
		x := l1.From.x()
		y := l2.y(x)
		return Latlong{y,x}, false

	} else if math.IsInf(b,0) {
		// l2 is vertical; as above, but switch the lines
		x := l2.From.x()
		y := l1.y(x)
		return Latlong{y,x}, false

	} else {
		x := (d-c) / (a-b)
		y := (a*d - b*c) / (a-b)
		return Latlong{y,x}, false
	}
}

// }}}

// {{{ l.PerpendicularTo

func (orig LatlongLine)PerpendicularTo(pos Latlong) LatlongLine {
	// The perpendicular has a gradient that is the negative inverse of the orig line
	m := -1 / orig.m
	perp := LatlongLine{
		From: pos,
		m: m,
		b: calcB(m, pos),
	}

	// Both lines have equations, and anchor points at l.From; we can intersect them to
	// derive the endpoint of the perpendicular. No chance of them being parallel :)
	perp.To,_ = orig.intersectByLineEquations(perp)

	return perp
}

// }}}
// {{{ l.ClosestTo

// Presumes infinite line
func (line LatlongLine)ClosestTo(pos Latlong) Latlong {
	perp := line.PerpendicularTo(pos) // end point of this is the intersection point
	return perp.To
}

// }}}
// {{{ l.ClosestDistance

func (line LatlongLine)ClosestDistance(pos Latlong) float64 {
	return pos.Dist(line.ClosestTo(pos))
}

// }}}
// {{{ l.DistAlongLine

// If one unit is the dist between .From and .To, and .From is zero; how far along the line is pos?
func (line LatlongLine)DistAlongLine(pos Latlong) float64 {
	// Todo: project pos onto the line itself, before doing the d1,d2,dPos stuff
	// That means a perpendicular; measure its length and discard if close to zero

	// The perpendicular will connect 'pos' to 'line'; the perp's endpoint will lie on 'line'
	// (although, latlong geometry is skew, so this is never really 'perpendicular' :/
	perp := line.PerpendicularTo(pos)
	
	// If line is more horizontal than vertical, project onto X axis; else Y
	d1,d2,dPos := 0.0,0.0,0.0
	if math.Abs(line.m) < 1.0 {
		d1,d2,dPos = line.From.x(),line.To.x(),perp.To.x()
	} else {
		d1,d2,dPos = line.From.y(),line.To.y(),perp.To.y()
	}

	// d1 represents 0.0; d2 represents 1.0. Where is pos ?
	return (dPos - d1) / (d2 - d1)
}

// }}}
// {{{ l.IntersectsUnbounded

// Treats the lines as infinite. Returns whether there was an intersection.
func (l1 LatlongLine)IntersectsUnbounded(l2 LatlongLine) (Latlong, bool) {
	pos,parallel := l1.intersectByLineEquations(l2)
	return pos, !parallel
}

// }}}
// {{{ l.Intersects

// Returns point of intersection (may be invalid), and bool stating if intersection occurred
func (l1 LatlongLine)Intersects(l2 LatlongLine) (Latlong, bool) {
	pos,parallel := l1.intersectByLineEquations(l2)

	if parallel { return pos, false }
	
	// Does the point of intersection lie within [from,to] for both lines ?
	// Simple bounding box tests will work !
	if ! l1.Box().Contains(pos) { return pos, false }
	if ! l2.Box().Contains(pos) { return pos, false }
	
	return pos, true
}

// }}}
// {{{ l.WhichSide

// -ve == left, +ve == right, 0 == lies-on-line
func (l LatlongLine)WhichSide(p Latlong) int {
	x,y := p.x(),p.y()
	x1,y1 := l.From.x(),l.From.y()
	x2,y2 := l.To.x(),l.To.y()

	d := (x - x1)*(y2 - y1) - (y - y1)*(x2 - x1)

	if d < 0.0 { return +1 }
	if d > 0.0 { return -1 }
	return 0 // Lies on the line
}

// }}}

// {{{ latlong.LiesOn

func (pos Latlong)LiesOn(line LatlongLine) bool {
	if ! line.From.BoxTo(line.To).Contains(pos) { return false }
	if line.ClosestDistance(pos) > kLineSnapKM  { return false }
	return true
}

// }}}

// {{{ -------------------------={ E N D }=----------------------------------

// Local variables:
// folded-file: t
// end:

// }}}

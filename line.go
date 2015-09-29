package geo

import(
	"fmt"
	"math"
)

// All of this works in simple (x,y) space. We take Long as x, to make horiz/vert look normal
func (ll Latlong)x() float64 { return ll.Long }
func (ll Latlong)y() float64 { return ll.Lat }

const kLineSnapKM = 0.3  // How far a trackpoint can be from a line, and still be on that line
const KLineSnapKM = 0.3  // How far a trackpoint can be from a line, and still be on that line

// Always using two anchor points, and then derive the equation of the line: y = m.x + b
type LatlongLine struct {
	From,To  Latlong
	m,b      float64
}
func (line LatlongLine)String() string {
	return fmt.Sprintf("[y=%.2f.x + %.2f] (%.1f,%.1f)->(%.1f.%.1f)",
		line.m, line.b, line.From.x(), line.From.y(), line.To.x(), line.To.y())
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

// {{{ latlong.BuildLine

func (from Latlong)BuildLine(to Latlong) LatlongLine {
	m := calcM(from,to)
	return LatlongLine{
		From: from,
		To: to,
		m: m,
		b: calcB(m,to),
	}
}

// }}}

// {{{ line.y

// Apply equation of line
func (line LatlongLine)y(x float64) float64 { return line.m * x + line.b }

// }}}
// {{{ line.buildPerpendicular

func (orig LatlongLine)buildPerpendicular(pos Latlong) LatlongLine {
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
// {{{ line.intersectByLineEqutions

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

// {{{ line.ClosestTo

// Presumes infinite line
func (line LatlongLine)ClosestTo(pos Latlong) Latlong {
	perp := line.buildPerpendicular(pos) // end point of this is the intersection point
	return perp.To
}

// }}}
// {{{ line.ClosestDistance

func (line LatlongLine)ClosestDistance(pos Latlong) float64 {
	return pos.Dist(line.ClosestTo(pos))
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

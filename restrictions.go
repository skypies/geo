package geo

import(
	"encoding/gob"
	"fmt"
)

func init() {
	gob.Register(new(DebugLog))
	gob.Register(SquareBoxRestriction{})
	gob.Register(VerticalPlaneRestriction{})
	gob.Register(PolygonRestriction{})
}

type DebugLog string
func (dl *DebugLog)GetDebug() string {
	ret := string(*dl)
	*dl = ""
	return ret
}
func (dl *DebugLog)Debugf(format string,args ...interface{}) {
	*dl += DebugLog(fmt.Sprintf(format,args...))
}


// SquareBoxRestricton fully implements geo.Restrictor
type SquareBoxRestriction struct {
	NamedLatlong             // embed
	SideKM                   float64
	AltitudeMin,AltitudeMax  int64
	IsExcluding              bool
	Debugger                 // embed; populate with ptr rcvr e.g. sb.Debugger = new(geo.DebugLog)
}
func (sb SquareBoxRestriction)String() string {
	nstr := sb.Name
	if nstr == "" { nstr = sb.Latlong.String() }
	str := fmt.Sprintf("SquareBox %s @%.1fKM", nstr, sb.SideKM)
	if sb.AltitudeMin > 0 || sb.AltitudeMax > 0 {
		str += fmt.Sprintf(" [%d,", sb.AltitudeMin)
		if sb.AltitudeMax>0 { str += fmt.Sprintf("%d",sb.AltitudeMax) } else { str += "-" }
		str += "]ft"
	}

	if sb.IsExcluding { str += "(EXCLUDES)" }
	return str
}

func (sb SquareBoxRestriction)IsExclusion() bool { return sb.IsExcluding }
func (sb SquareBoxRestriction)IsNil() bool { return sb.SideKM==0 || sb.NamedLatlong.IsNil() }

func (sb SquareBoxRestriction)ToCircles() []LatlongCircle { return nil }
func (sb SquareBoxRestriction)ToLines() []LatlongLine {
	return sb.Box(sb.SideKM,sb.SideKM).ToLines()
}

func (sb SquareBoxRestriction)BoundingBox() LatlongBox { return sb.Box(sb.SideKM,sb.SideKM) }
func (sb SquareBoxRestriction)CanContain() bool { return true }
func (sb SquareBoxRestriction)Contains(pos Latlong) bool {
	return sb.Box(sb.SideKM,sb.SideKM).Contains(pos)
}
func (sb SquareBoxRestriction)OverlapsLine(ln LatlongLine) OverlapOutcome {
	return sb.Box(sb.SideKM,sb.SideKM).OverlapsLine(ln)
}
func (sb SquareBoxRestriction)OverlapsAltitude(a int64) OverlapOutcome {
	// r2 is the altitude; so if too low, it comes 'before' the restriction
	if sb.AltitudeMin > 0 &&  a < sb.AltitudeMin { return DisjointR2ComesBefore }
	if sb.AltitudeMax > 0 &&  a > sb.AltitudeMax { return DisjointR2ComesAfter }
	return OverlapR2IsContained
}

// VerticalPlane  fully implements geo.Restrictor
type VerticalPlaneRestriction struct {
	Start, End               NamedLatlong
	AltitudeMin,AltitudeMax  int64
	IsExcluding              bool
	Debugger                 // embed; populate with ptr rcvr e.g. vp.Debugger = new(geo.DebugLog)
}
func (vp VerticalPlaneRestriction)String() string {
	nstr1,nstr2 := vp.Start.Name,vp.End.Name
	if nstr1 == "" { nstr1 = vp.Start.String() }
	if nstr2 == "" { nstr2 = vp.End.String() }
	str := fmt.Sprintf("VerticalPlane %s to %s", nstr1, nstr2)
	if vp.AltitudeMin > 0 || vp.AltitudeMax > 0 {
		str += fmt.Sprintf(" [%d,", vp.AltitudeMin)
		if vp.AltitudeMax>0 { str += fmt.Sprintf("%d",vp.AltitudeMax) } else { str += "-" }
		str += "]ft"
	}
	if vp.IsExcluding { str += "(EXCLUDES)" }
	return str
}
func (vp VerticalPlaneRestriction)IsExclusion() bool { return vp.IsExcluding }
func (vp VerticalPlaneRestriction)IsNil() bool { return vp.Start.IsNil() || vp.End.IsNil() }

func (vp VerticalPlaneRestriction)ToCircles() []LatlongCircle { return nil }
func (vp VerticalPlaneRestriction)ToLines() []LatlongLine {
	return []LatlongLine{vp.Start.LineTo(vp.End.Latlong)}
}

func (vp VerticalPlaneRestriction)BoundingBox() LatlongBox { return vp.Start.BoxTo(vp.End.Latlong) }
func (vp VerticalPlaneRestriction)CanContain() bool { return false }
func (vp VerticalPlaneRestriction)Contains(pos Latlong) bool { return false }

func (vp VerticalPlaneRestriction)OverlapsLine(ln LatlongLine) OverlapOutcome {
	if _,intersects := vp.Start.LineTo(vp.End.Latlong).Intersects(ln); intersects {
		// The only meaningful way to express "ln intersects the vp" as an OverlapOutcome
		return OverlapR2StraddlesStart
	} else {
		return Disjoint
	}
}
func (vp VerticalPlaneRestriction)OverlapsAltitude(a int64) OverlapOutcome {
	// r2 is the altitude; so if too low, it comes 'before' the restriction
	if vp.AltitudeMin > 0 &&  a < vp.AltitudeMin { return DisjointR2ComesBefore }
	if vp.AltitudeMax > 0 &&  a > vp.AltitudeMax { return DisjointR2ComesAfter }
	return OverlapR2IsContained
}


// PolygonRestricton fully implements geo.Restrictor
type PolygonRestriction struct {
	*Polygon
	SideKM                   float64
	AltitudeMin,AltitudeMax  int64
	IsExcluding              bool
	Debugger                 // embed; populate with ptr rcvr e.g. pr.Debugger = new(geo.DebugLog)
}
func (pr PolygonRestriction)String() string {
	str := fmt.Sprintf("%d-gon ~%.2fKM @ %s", len(pr.Polygon.Path.Points()),
		pr.Polygon.ApproxRadiusKM(), pr.Polygon.Centroid())

	if pr.AltitudeMin > 0 || pr.AltitudeMax > 0 {
		str += fmt.Sprintf(" [%d,", pr.AltitudeMin)
		if pr.AltitudeMax>0 { str += fmt.Sprintf("%d",pr.AltitudeMax) } else { str += "-" }
		str += "]ft"
	}

	if pr.IsExcluding { str += "(EXCLUDES)" }
	return str
}

func (pr PolygonRestriction)IsExclusion() bool { return pr.IsExcluding }
func (pr PolygonRestriction)IsNil() bool {
	return pr.Polygon == nil || len(pr.Polygon.Path.Points()) < 3
}

func (pr PolygonRestriction)ToCircles() []LatlongCircle { return nil }
func (pr PolygonRestriction)ToLines() []LatlongLine { return pr.Polygon.ToLines() }

func (pr PolygonRestriction)BoundingBox() LatlongBox { return LatlongBoxFromBound(pr.Bound()) }
func (pr PolygonRestriction)CanContain() bool { return true }
//func (pr PolygonRestriction)Contains(pos Latlong) bool { return pr.Polygon.Contains(pos) }
//func (pr PolygonRestriction)OverlapsLine(ln LatlongLine) OverlapOutcome {
//	return pr.Box(pr.SideKM,pr.SideKM).OverlapsLine(ln)
//}
func (pr PolygonRestriction)OverlapsAltitude(a int64) OverlapOutcome {
	// r2 is the altitude; so if too low, it comes 'before' the restriction
	if pr.AltitudeMin > 0 &&  a < pr.AltitudeMin { return DisjointR2ComesBefore }
	if pr.AltitudeMax > 0 &&  a > pr.AltitudeMax { return DisjointR2ComesAfter }
	return OverlapR2IsContained
}

package geo

// go test -v github.com/skypies/geo

import "fmt"
import "testing"

func TestPolygon(t *testing.T) {
	poly := NewPolygon()
	poly.AddPoint(Latlong{ 0, 0})
	poly.AddPoint(Latlong{ 0,10})
	poly.AddPoint(Latlong{10,10})
	poly.AddPoint(Latlong{10, 0})

	if !poly.Centroid().Equal(Latlong{5,5}) {
		t.Errorf("Centroid not 5,5: %s\n", poly)
	}

	intersectionTests := []struct{
		A,B Latlong
		N   int
	}{
		{Latlong{  0, 20}, Latlong{ 10, 20}, 0},
		{Latlong{  5,  5}, Latlong{ 15,  5}, 1},
		{Latlong{ -5,  5}, Latlong{ 15,  5}, 2},
		// Corner tests
		{Latlong{  0,  0}, Latlong{ 10, 10}, 2},
		{Latlong{ 10, 10}, Latlong{ 10, 20}, 2}, // We get two, as the line is collinear with poly
	}

	for i,test := range intersectionTests {
		l := test.A.LineTo(test.B)
		pts,plcs := poly.IntersectsLine(l)
		if len(pts) != test.N {
			for _,pt := range pts {
				fmt.Printf(" ** %#v, %#v\n", pt,plcs)
			}
			t.Errorf("IntersectsLine[%3d]: expected %d, saw %d. L:%s\n", i, test.N, len(pts), l)
		}
	}
}	

/*
func TestContains(t *testing.T) {
	poly := NewPolygon()
	poly.AddPoint(Latlong{  0,  0})
	poly.AddPoint(Latlong{  0, 10})
	poly.AddPoint(Latlong{  5,  5}) // Instead of a square, dip into center like an "M"
	poly.AddPoint(Latlong{ 10, 10})
	poly.AddPoint(Latlong{ 10,  0})

	containsTests := []struct{
		Expected bool
		A Latlong
	}{
		{false, Latlong{  0, 20}},
		{true,  Latlong{  5,  2}},
		// The concave void
		{false, Latlong{  5,  5.1}}, // just outside
		{true , Latlong{  5,  4.9}}, // just inside
		// Corners; should all be true, but the intersection code isn't robust
		{false, Latlong{  5,  5}}, // Shoud be true
		{true,  Latlong{ 10, 10}},
		{true,  Latlong{  0, 10}},
	}

	for i,test := range containsTests {
		actual := poly.Contains(test.A)
		if actual != test.Expected {
			t.Errorf("Contains[%3d]: expected %v, saw %v. Point:%s\n", i, test.Expected, actual, test.A)
		}
	}
}
*/

func TestOverlapsLine(t *testing.T) {
	poly := NewPolygon()
	poly.AddPoint(Latlong{  0,  0})
	poly.AddPoint(Latlong{  0, 10})
	poly.AddPoint(Latlong{  5,  5}) // Instead of a square, dip into center like an "M"
	poly.AddPoint(Latlong{ 10, 10})
	poly.AddPoint(Latlong{ 10,  0})

	tests := []struct{
		Expected OverlapOutcome
		A,B Latlong
	}{
		{Disjoint,                Latlong{-10, 50},   Latlong{10,50}},
		{Disjoint,                Latlong{  5,  5.1}, Latlong{ 5, 6}}, // in the cavity
		{OverlapR2Contains,       Latlong{-10,  2},   Latlong{20, 2}},
		{OverlapR2IsContained,    Latlong{  2,  2},   Latlong{ 3, 3}},
		{OverlapR2StraddlesStart, Latlong{ -4,  4},   Latlong{ 3, 3}},
		{OverlapR2StraddlesEnd,   Latlong{  4,  4},   Latlong{ 3, 13}},
	}

	for i,test := range tests {
		actual := poly.OverlapsLine(test.A.LineTo(test.B))
		if actual != test.Expected {
			t.Errorf("OverlapsLine[%3d]: expected %v, saw %v. Line:%s - %s\n", i, test.Expected, actual, test.A, test.B)
		}
	}
}

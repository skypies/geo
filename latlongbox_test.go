package geo
// go test -v github.com/skypies/geo

import(
	"testing"
)

func TestBoxOverlapsLine(t *testing.T) {
	box := Latlong{0,0}.BoxTo(Latlong{10,10})

	tests := []struct{
		Expected OverlapOutcome
		S,E Latlong
	}{
		// Basic tests
		{Disjoint,                Latlong{15,15}, Latlong{16,16}},
		{OverlapR2IsContained,    Latlong{ 2, 2}, Latlong{ 3, 3}},

		// R1 is the box; R2 is the line. "line straddles start of box" means line starts outside.
		{OverlapR2StraddlesStart, Latlong{20, 2}, Latlong{ 2, 2}},   // line starts outside
		{OverlapR2StraddlesEnd,   Latlong{ 2, 2}, Latlong{20, 2}},   // line starts inside

		// Corner tests
		{Disjoint,                Latlong{ 0,20}, Latlong{20, 0.1}}, // bbox overlap, but misses
		{OverlapR2Contains,       Latlong{ 0,20}, Latlong{20, 0}},   // line touches corner
		{OverlapR2StraddlesEnd,   Latlong{ 10,10},Latlong{12,12}},   // line starts at corner
		{OverlapR2StraddlesStart, Latlong{ 20,20},Latlong{10,10}},   // line ends at corner
	}

	for i,test := range tests {
		actual := box.OverlapsLine(test.S.LineTo(test.E))
		if actual != test.Expected {
			t.Errorf("[test % 2d] IntersectsLine: expected %v, got %v", i, test.Expected, actual)
		}
	}
}

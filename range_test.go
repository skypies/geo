package geo

import "testing"

type FloatTest struct {
	u,v, r,s  float64
	o         OverlapOutcome
}

func TestFloatOverlaps(t *testing.T) {
	tests := []FloatTest{
		{ 1.0, 2.0,    3.0, 4.0,    DisjointR2ComesAfter},
		{10.0,12.0,    3.0, 4.0,    DisjointR2ComesBefore},
		{ 1.0, 2.0,    0.5, 1.5,    OverlapR2StraddlesStart},
		{ 1.0, 2.0,    1.5, 2.5,    OverlapR2StraddlesEnd},
		{ 1.0, 2.0,    1.2, 1.8,    OverlapR2IsContained},
		{ 1.0, 2.0,    0.5, 4.0,    OverlapR2Contains},

		{ 1.0, 2.0,    1.0, 2.0,    OverlapR2IsContained},
	}
	for i,test := range tests {
		r1,r2 := Float64Range{test.u, test.v}, Float64Range{test.r, test.s}
		actual := RangeOverlap(r1,r2)
		if actual != test.o {
			t.Errorf("[t%d] exepected %v, got %v [%f,%f]", i, test.o, actual,
			)
		}
	}
}

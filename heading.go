package geo

// A few crappy routines for handling arithmetic on headings

// This float is supposed to range from [0.0,wrap) - i.e. a heading, [0,360)
// Pick the shortest path across the boundary (i.e. 6->354 is a delta of 12, not 348)
func interpolateWrappingFloat64(from, to, ratio, wrap float64) float64 {
	deltaCW  := to-from            // clockwise:         354 - 6     = 348
	deltaCCW := (from+wrap) - to   // counter clockwise: (6+360)-354 = 12

	if deltaCW <= deltaCCW {
		// We should head clockwise - just like interpolatefLoat64
		return from + (deltaCW)*ratio
	} else {
		// Take a shortcut over the wrap boundary
		new := from - (deltaCCW)*ratio
		if new < 0.0 { new += wrap }
		return new
	}
}

func InterpolateHeading(from, to, ratio float64) float64 {
	return interpolateWrappingFloat64(from, to, ratio, 360.0)
}

func HeadingDelta(from, to float64) float64 {
	delta := to - from
	if delta < -180.0 { delta += 360.0 }
	if delta >= 180.0 { delta -= 360.0 }
	return delta
}

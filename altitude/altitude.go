package altitude

// Aviation altitudes are complex. Read about them on Wikipedia.

// Given a pressure altitude, and a current value for the local atmospheric pressure,
// blah blah

const InHgAtStandardAtmosphericPressure = 29.9213

// When pressure is higher than normal, the pressure sensor in the altimeter will register
// a higher pressure - which means lower down in the atmosphere - so the pressure altitude
// will be too low. So we must add the appropriate amount of height.
func PressureAltitudeToIndicatedAltitude(p, inHg float64) float64 {
	// Implement for realz
	deltaInMg := inHg - InHgAtStandardAtmosphericPressure

	return p + (96.0 * (deltaInMg/0.1))  // cf. {100' for every 0.1" of Hg}
}


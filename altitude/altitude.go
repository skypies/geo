package altitude
// Aviation altitudes are complex. Read about them on Wikipedia.
// https://en.wikipedia.org/wiki/Altitude#Altitude_in_aviation

import "math"

const (
	// Some universal physical constants, using american units
	Rstar = 8.9494596e4    // Universal gas constant    (lb.ft2/(lb.mol.s2))
	Gzero = 32.17405       // Gravitational constant    (ft/s^2)
	M     = 28.9644        // Molar mass of Earth's air (lb/lb.mol)

	// These are for the 'b=0' regime, which applies under 36,000 feet.
	p0    = 29.9213           // Static pressure        (inHg)
	T0    = 288.15            // Standard Temperature   (K)
	L0    = -0.0019812        // Temperature lapse rate (K/ft)


	// Precompute these values
	E              = (Gzero * M) / (Rstar * L0)  // This is 'E' in the derivation below
	OneOverMinusE  = -1.0/E                      // The final exponent
	T0OverMinusL0  = -1.0 * (T0/L0)
)

// The barometric formula generates a 'pressure altitude' (AKA
// 'Standard Height'), when given a pressure reading from an
// altimeter. This is based on the Standard Atmosphere; as pressure
// readings drop, the computed pressure-altitude increases.

// In real atmospheric conditions, the local pressure may be higher or
// lower; this causes pressure altitudes to differ significantly from
// the actual height from sea level.

// E.g. if the local pressure is lower than normal, the barometric
// formula will interpret that as being higher up in the atmosphere;
// so the pressure-altitude value will be too high, as the aircraft
// will physically be lower.

// Given the local pressure reading, we use the barometric formula to
// compute the corresponding height in the 'standard atmosphere'. We
// apply this as on offset against the supplied pressure-altitude, and
// generate an altitude that is much closer to the physical height
// from the ground - the 'indicated altitude'.

func PressureAltitudeToIndicatedAltitude(pressureAlt, inHg float64) float64 {
	return pressureAlt - StandardHeightFromBarometricPressure(inHg)
}

func StandardHeightFromBarometricPressure(bp float64) float64 {
	//   h = T/-L . ( 1 - (BP/P)^(1/-E) )
	return T0OverMinusL0 * (1.0 - math.Pow(bp/p0, OneOverMinusE))
}

/* Derivation

Start with the barometric formula, and the constants for the US standard atmosphere.
https://en.wikipedia.org/wiki/Barometric_formula

      BP = Pb . (Tb / (Tb + Lb(h-hb))) ^ (g0.M / R*.Lb)

BP is the barometric pressure we have a measurement of. (I've renamed
it BP to avoid confusion with P.)

We're only concerned with altitudes under 36,000' so we're in regime
b=0, and so:

  hb := 0

For clarity, I'll omit the 'b' subscripts and name the constant exponent
value as E:

  E := g0.M / R*.Lb

Thus the barometric formula becomes:

                    BP = P . ( T / (T+Lh) )^E

We reformulate it in terms of h (height):

         BP          = P . ( T / (T+Lh) )^E
         BP/P        =     ( T / (T+Lh) )^E
        (BP/P)^(1/E) =       T / (T+Lh)                    [take the E'th root of both sides]
 (T+Lh).(BP/P)^(1/E) =       T
  T+Lh               =       T . (P/BP)^(1/E)
    Lh               =       T . (P/BP)^(1/E) - T
    Lh               =       T . ( (P/BP)^(1/E) - 1 )
     h               =     T/L . ( (P/BP)^(1/E) - 1 )

We can further tidy it up; the constants result in L and EXP being
negative value, so we shuffle around to flip a few signs:

     h               =     T/L . ( (BP/P)^(1/-E) - 1 )     [using x^y == 1/(x^-y)]
     h               =    T/-L . ( 1 - (BP/P)^(1/-E) )     [mult by -1/-1]

The values for Tb, hb and Lb for the regime where b=0:

  T := 288.15 K
  h := 0 ft
  L := -0.0019812 K/ft

and the fixed constants:

  P0 := 29.92126 inHg
  R* := 8.9494596.10^4
  g0 := 32.17405
  M  := 28.9644

So:

     E := (32.17405 * 28.9644) / (8.9494596.10^4 * -0.0019812)
        = -5.2558763
  1/-E := 0.190263 

  h = 288.15/0.0019812 . ( 1 - (BP/29.92126)^0.190263 )
  h = 145442 . ( 1 - (BP/29.92126)^0.190263 )

Sanity check: a pilot's rule of thumb is "every 0.1" of mercury is
about 100 feet of altitude". So we set the observed barometric
pressure to be 0.1" less than standard pressure ...

  BP := P0       - 0.1
      = 29.92126 - 0.1
      = 29.82126

... and drop it into the simplifed formula to see what height
corresponds to the pressure reading ...

  h = 145442 . ( 1 - (BP/29.92126)^0.190263 )
    = 145442 . ( 1 - (29.82126/29.92126)^0.190263 )
    = 92.6

... that's close enough to 100' !

*/

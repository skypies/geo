package geo

// NamedLatlong is a Latlong, that happens to have a name (typically a waypoint)
type NamedLatlong struct {
	Name     string
	Latlong  // embedded
}

func (nl NamedLatlong)String() string {
	return nl.Latlong.String() + "[" + nl.Name + "]"
}
func (nl NamedLatlong)ShortString() string {
	if nl.Name != "" {
		return nl.Name
	}
	return nl.Latlong.String()
}

func (nl NamedLatlong)IsNil() bool {
	return nl.Name == "" && nl.Latlong.IsNil()
}

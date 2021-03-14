package arb_converter

type DataArb struct {
	Cultures []string
	Items    map[string]*ItemArb
}

type ItemArb struct {
	Description string
	Cultures    map[string]string
	Parameters  map[string]struct{}
}

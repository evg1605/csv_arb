package arb

type Data struct {
	Cultures []string
	Items    map[string]*Item
}

type Item struct {
	Description string
	Cultures    map[string]string
	Parameters  map[string]struct{}
}

package sky

// Link link
type Link struct {
	Href  string
	Label string
}

// Dropdown dropdown
type Dropdown struct {
	Label string
	Links []*Link
}

package sky

const (
	// LayoutApplication application layout
	LayoutApplication = "application"
	// LayoutDashboard dashboard layout
	LayoutDashboard = "dashboard"
)

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

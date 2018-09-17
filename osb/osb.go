package osb

// Catalog represents object of OpenServiceBroker API
type Catalog struct {
	Services []*Service `json:"services"`
}

// Service represents object of OpenServiceBroker API
type Service struct {
	Name        string   `json:"name"`
	ID          string   `json:"id"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Requires    []string `json:"requires"`
	Bindable    bool     `json:"bindable"` // nolint
	// Metadata        *Metadata        `json:"metadata,omitempty"`
	PlanUpdateable bool    `json:"plan_updateable,omitempty"`
	Plans          []*Plan `json:"plans"`
}

// Plan represents object of OpenServiceBroker API
type Plan struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	// Metadata    *Metadata `json:"metadata,omitempty"`
	Bindable bool `json:"bindable"`
	Free     bool `json:"free"`
}

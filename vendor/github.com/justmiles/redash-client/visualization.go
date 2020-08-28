package redash

// Visualization redash viz
type Visualization struct {
	Description string `json:"description,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
	ID          int    `json:"id,omitempty"`
	Type        string `json:"type,omitempty"`
	Name        string `json:"name,omitempty"`
}

// Visualizations is a slice of Visualization
type Visualizations []Visualization

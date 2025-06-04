package health

const (
	UP   = "UP"
	DOWN = "DOWN"
)

type Health struct {
	Status  string         `json:"status,omitempty"`
	Details map[string]any `json:"details,omitempty"`
}

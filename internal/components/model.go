package components

// Component represents a stored React component.
type Component struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Code        string            `json:"code"`
	PropsSchema map[string]string `json:"props_schema"`
}

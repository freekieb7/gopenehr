package webhook

type CloudEvent struct {
	SpecVersion string         `json:"specversion"`
	Type        string         `json:"type"`
	Source      string         `json:"source"`
	ID          string         `json:"id"`
	Time        string         `json:"time"`
	Data        map[string]any `json:"data"`
}

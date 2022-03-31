package types

type Step struct {
	Type        string `json:"type"`
	Source      string `json:"source"`
	Destination string `json:"dest"`
	Embed       string `json:"embed"`
}

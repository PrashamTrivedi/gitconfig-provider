package cmd

type GitProvider struct {
	Name       string            `json:"name,omitempty"`
	Url        string            `json:"url,omitempty"`
	Properties map[string]string `json:"properties,omitempty"`
}

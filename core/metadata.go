package space

type Metadata struct {
	Title      string                 `yaml:"title"`
	Date       string                 `yaml:"date"`
	Layout     string                 `yaml:"layout"`
	Permalink  string                 `yaml:"permalink"`
	Categories []string               `yaml:"categories"`
	Tags       []string               `yaml:"tags"`
	Drafts     bool                   `yaml:"drafts"`
	Vars       map[string]interface{} `yaml:"vars,omitempty"`
}

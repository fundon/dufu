package space

type Metadata struct {
	Title      string
	Date       string
	Layout     string
	Permalink  string
	Categories []string
	Tags       []string
	Drafts     bool
	Others     map[string]interface{} `yaml:"others,omitempty`
}

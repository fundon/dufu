package space

type Site struct {
	Author      Author
	Pages       []Page
	Posts       []Page
	Tags        []string
	Categories  []string
	Url         string
	Source      string
	Destination string
	Time        string
}

type Author struct {
	Name     string
	FullName string
	Emial    string
}

type Page struct {
	Title      string
	Url        string
	Date       string
	Id         string
	Path       string
	Content    string
	Tags       []string
	Categories []string
	Next       *Page
	Previous   *Page
}

type Metadata struct {
	Title      string   `yaml:"title"`
	Date       string   `yaml:"date"`
	Layout     string   `yaml:"layout"`
	Permalink  string   `yaml:"permalink"`
	Draft      bool     `yaml:"draft"`
	Categories []string `yaml:"categories"`
	Tags       []string `yaml:"tags"`
	Vars       Locals   `yaml:"vars,omitempty"`
}

type Locals map[string]interface{}

package meta

type Metadata struct {
	Title       string   `yaml:"title"`
	Date        string   `yaml:"date"`
	Slug        string   `yaml:"slug"`
	Tags        []string `yaml:"tags"`
	Description string   `yaml:"description"`
	Draft       bool     `yaml:"draft"`
}

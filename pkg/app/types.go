package app

type Config struct {
	Units []PushUnit `yaml:"units"`
}

type PushUnit struct {
	Id             string `yaml:"id"`
	Repository     string `yaml:"repository"`
	RepositoryPath string
	Files          []string `yaml:"files"`
	Exclude        []string `yaml:"exclude"`
	CommitMessage  string   `yaml:"commit_message"`
	AuthorName     string   `yaml:"author_name"`
	AuthorEmail    string   `yaml:"author_email"`
	SshKey         string   `yaml:"ssh_key"`
}

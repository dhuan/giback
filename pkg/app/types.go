package app

type Config struct {
	Units []PushUnit
}

type PushUnit struct {
	Id             string
	Repository     string
	RepositoryPath string
	Files          []string
	Exclude        []string
	Commit_Message string
	Author_Name    string
	Author_Email   string
	Ssh_Key        string
}

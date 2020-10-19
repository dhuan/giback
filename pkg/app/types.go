package app

type Config struct {
	Units []PushUnit
}

type PushUnit struct {
	Id         string
	Repository string
	Key        string
	Files      []string
}

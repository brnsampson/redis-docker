package instance

type InstanceInfoParser interface {
	ParseInstanceInfo(string, *InstanceInfo) error
}

type InstanceShimmer interface {
	Ping() error
	GetInfo() (string, error)
	ReadConfig(string) (string, error)
	UpdateConfig(string, string) error
	SetReplication(string) error
	Quit() error
}

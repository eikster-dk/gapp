package secrets

type Spinner interface {
	Start() error
	Message(msg string)
	Stop() error
	Fail() error
}

type Secret struct {
	Name  string
	Value string
}

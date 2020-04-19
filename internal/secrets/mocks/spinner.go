package mocks

type NoopSpinner struct{}

func (f NoopSpinner) Start() error {
	return nil
}

func (f NoopSpinner) Message(msg string) {
}

func (f NoopSpinner) Stop() error {
	return nil
}

func (f NoopSpinner) Fail() error {
	return nil
}

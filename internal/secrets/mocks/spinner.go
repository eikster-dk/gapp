package mocks

type MockSpinner struct{}

func (m MockSpinner) Start() error {
	return nil
}

func (m MockSpinner) Message(msg string) {
}

func (m MockSpinner) Stop() error {
	return nil
}

func (m MockSpinner) Fail() error {
	return nil
}

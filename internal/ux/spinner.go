package ux

import (
	"github.com/theckman/yacspin"
	"time"
)

type Spinner struct {
	spinner *yacspin.Spinner
}

func NewSpinner(suffix, message string) (*Spinner, error) {
	cfg := yacspin.Config{
		Frequency:       200 * time.Millisecond,
		CharSet:         yacspin.CharSets[59],
		Suffix:          suffix,
		SuffixAutoColon: true,
		Message:         message,
		StopCharacter:   "✓",
		StopColors:      []string{"fgGreen"},
	}

	spinner, err := yacspin.New(cfg)
	if err != nil {
		return nil, err
	}

	return &Spinner{
		spinner: spinner,
	}, nil
}

func (s *Spinner) Start() error {
	return s.spinner.Start()
}

func (s *Spinner) Message(msg string) {
	s.spinner.Message(msg)
}

func (s *Spinner) Stop() error {
	return s.spinner.Stop()
}

func (s *Spinner) Fail() error {
	return s.spinner.StopFail()
}

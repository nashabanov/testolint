package no_suites

import "github.com/ozontech/testo"

var _ *testo.T

type Ordinary struct{}

func (Ordinary) TestBroken()     {}
func (Ordinary) BeforeEach() int { return 1 }

package hooksignature

import "github.com/ozontech/testo"

type T struct{ *testo.T }

type Suite struct { // want `TESTO011: suite "Suite" contains no tests`
	testo.Suite[T]
}

func (Suite) BeforeAll(t T) {}

func (Suite) BeforeEach() {} // want "TESTO002"

func (Suite) AfterEach(t string) {} // want "TESTO002"

func (Suite) AfterAll(t T) error { // want "TESTO002"
	return nil
}

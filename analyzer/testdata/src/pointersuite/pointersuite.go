package pointersuite

import "github.com/ozontech/testo"

type T struct{}

type Suite struct {
	*testo.Suite[T]
}

func (Suite) TestFoo() {} // want "suite Suite: found test TestFoo"

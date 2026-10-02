package aliasimport

import t "github.com/ozontech/testo"

type T struct{}

type Suite struct {
	t.Suite[T]
}

func (Suite) TestFoo() {} // want "suite Suite: found test TestFoo"

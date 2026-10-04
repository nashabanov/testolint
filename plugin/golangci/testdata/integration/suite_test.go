package integration

import "github.com/ozontech/testo"

type Suite struct{ testo.Suite[*testo.T] }

// Deliberately missing the Testo T argument: TESTO001.
func (Suite) TestBroken() {}

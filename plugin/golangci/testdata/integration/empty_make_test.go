package integration

import "github.com/ozontech/testo"

type EmptyMakeSuite struct{ testo.Suite[*testo.T] }

func (EmptyMakeSuite) TestValue(t *testo.T, p struct{ Value int }) {}
func (EmptyMakeSuite) CasesValue() []int                           { return make([]int, 0, 4) }

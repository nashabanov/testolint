package emptymakeshadow

import "github.com/ozontech/testo"

type Suite struct{ testo.Suite[*testo.T] }

func (Suite) TestValue(t *testo.T, p struct{ Value int }) {}
func (Suite) CasesValue() []int                           { return make([]int{}, 0) }
func make([]int, int) []int                               { return []int{1} }

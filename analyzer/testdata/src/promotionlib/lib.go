package promotionlib

import "github.com/ozontech/testo"

type Base struct{ testo.Suite[*testo.T] }

func (Base) CasesRole(x int) []string { return []string{"admin"} }

type Methods[T any] struct{}

func (Methods[T]) TestValue(t T, p struct{ Value T }) {}
func (*Methods[T]) CasesValue() []T                   { return []T{*new(T)} }

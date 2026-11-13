package lib

import "github.com/ozontech/testo"

type Base[T any] struct{ testo.Suite[T] }

func (Base[T]) BeforeEach(t T) {}

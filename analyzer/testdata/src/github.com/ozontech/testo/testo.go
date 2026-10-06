package testo

type T struct{}

type Suite[T any] struct{}

func (Suite[T]) private() {}

func (Suite[T]) BeforeAll(t T)  {}
func (Suite[T]) BeforeEach(t T) {}
func (Suite[T]) AfterEach(t T)  {}
func (Suite[T]) AfterAll(t T)   {}

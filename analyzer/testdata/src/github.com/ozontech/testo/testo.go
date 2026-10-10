package testo

type T struct{}

func (*T) Skip(args ...any) {}

type Suite[T any] struct{}

func (Suite[T]) private() {}

func (Suite[T]) BeforeAll(t T)  {}
func (Suite[T]) BeforeEach(t T) {}
func (Suite[T]) AfterEach(t T)  {}
func (Suite[T]) AfterAll(t T)   {}

// These constraints and type-parameter orders mirror Testo v1.8.0.
// CommonT and options are reduced to the surface needed by analyzer fixtures.
type CommonT interface{ privateT() }

func (*T) privateT() {}

type suite[T any] interface {
	private()
	BeforeAll(T)
	BeforeEach(T)
	AfterEach(T)
	AfterAll(T)
}

func RunSuite[S suite[T], T CommonT](t any, s S, options ...any) bool                   { return true }
func RunSubSuite[S suite[Sub], Parent, Sub CommonT](t Parent, s S, options ...any) bool { return true }

package basic

import "github.com/ozontech/testo"

type T struct{}

// Валидный Testo suite.
type FirstSuite struct {
	testo.Suite[T]
}

func (FirstSuite) TestFoo() {} // want "suite FirstSuite: found test TestFoo"

func (*FirstSuite) TestBar() {} // want "suite FirstSuite: found test TestBar"

func (FirstSuite) CasesFoo() []int {
	return nil
}

func (*FirstSuite) BeforeEach() {}

func (FirstSuite) AfterAll() {}

func (FirstSuite) Helper() {}

// Обычный тип — не должен считаться suite.
type NotTesto struct{}

func (NotTesto) TestShouldBeIgnored() {}

// Обычная package-level test function — тоже игнорируем.
func TestRegular() {}

// Testo suite без TestXxx.
type EmptySuite struct {
	testo.Suite[T]
}

func (EmptySuite) Helper() {}

// Ещё один Testo suite, чтобы проверить группировку типов.
type SecondSuite struct {
	testo.Suite[T]
}

func (SecondSuite) TestBaz() {} // want "suite SecondSuite: found test TestBaz"

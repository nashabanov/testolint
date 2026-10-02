package wrongpackage

import fake "faketesto"

type T struct{}

type Suite struct {
	fake.Suite[T]
}

func (Suite) TestFoo() {}

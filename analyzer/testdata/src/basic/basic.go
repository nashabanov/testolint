package basic

type Suite struct{}

func (Suite) TestFoo() {} // want "testolint: found new test method TestFoo"

func (*Suite) TestBar() {} // want "testolint: found new test method TestBar"

func (Suite) Foo() {}

func TestBaz() {}

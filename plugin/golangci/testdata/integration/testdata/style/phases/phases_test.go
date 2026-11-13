package phases

import (
	"github.com/ozontech/testo"
	allure "github.com/ozontech/testo-allure"
)

type T struct {
	*testo.T
	*allure.PluginAllure
}
type Suite struct{ testo.Suite[T] }

func action() {}
func verify() {}
func (Suite) TestOrder(t T) {
	action()
	verify()
	action()
	allure.Setup(t, "late setup", func(t T) {})
}

func (Suite) TestNativeAssertions(t T) {
	t.Require().Equal(1, 1)
	action()
	checks := t.Assert().Masked()
	checks.Equal(1, 1)
	action()
}

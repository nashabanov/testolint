package header

import (
	"github.com/ozontech/testo"
	allure "github.com/ozontech/testo-allure"
)

type T struct {
	*testo.T
	*allure.PluginAllure
}
type Suite struct{ testo.Suite[T] }

func (Suite) TestMissing(t T) {}
func (Suite) TestValid(t T)   { t.Title("title") }

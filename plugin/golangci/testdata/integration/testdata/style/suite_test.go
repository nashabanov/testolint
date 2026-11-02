package style

import (
	"testing"

	"github.com/ozontech/testo"
	allure "github.com/ozontech/testo-allure"
)

type T struct {
	*testo.T
	*allure.PluginAllure
}

type OrderSuite struct{ testo.Suite[T] }

func (OrderSuite) TestCreateOrder(t T) {
	t.Title("Create order")
	t.Tags("smoke")
	allure.Step(t, "create", func(T) {})
	allure.Step(t, "get", func(T) {})
	allure.Step(t, "check", func(T) {})
}

func TestOrder(t *testing.T) {
	t.Chdir(t.TempDir())
	testo.RunSuite(t, new(OrderSuite))
}

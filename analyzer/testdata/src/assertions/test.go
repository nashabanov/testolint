package assertions

import (
	"github.com/ozontech/testo"
	a "github.com/ozontech/testo-allure"
	check "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Context struct{ *a.PluginAllure }
type T = *Context
type Suite struct{ testo.Suite[T] }

func (*Suite) TestAssertions(t T) {
	check.Equal(t, 1, 1)   // want "TESTO031"
	require.Equal(t, 1, 1) // want "TESTO031"
	checks := check.New(t)
	checks.Equal(1, 1)         // want "TESTO031"
	require.New(t).Equal(1, 1) // want "TESTO031"
	_ = check.ObjectsAreEqual(1, 1)
	a.Step(t, "check", func(t T) { check.Equal(t, 1, 1) })
	a.Setup(t, "precondition", func(t T) { require.Equal(t, 1, 1) })
}

func returning(t T) bool { return check.Equal(t, 1, 1) }
func (*Suite) TestReturningHelper(t T) {
	_ = returning(t) // want "TESTO031"
}

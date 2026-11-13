package assertions

import (
	"github.com/ozontech/testo"
	a "github.com/ozontech/testo-allure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type T struct {
	*testo.T
	*a.PluginAllure
}
type Suite struct{ testo.Suite[T] }

func (Suite) TestAssertions(t T) {
	assert.Equal(t, 1, 1)
	require.New(t).Equal(1, 1)
	_ = assert.ObjectsAreEqual(1, 1)
	a.Setup(t, "precondition", func(t T) { require.Equal(t, 1, 1) })
}

package integration

import (
	"github.com/ozontech/testo"
	"testing"
)

func executionFailures(t *testing.T, ctx *testo.T) {
	testo.RunTest[*testo.T](t, nil)
	t.Run("nil", testo.Test[*testo.T](nil))
	ctx.Cleanup(func() { testo.Run(ctx, "cleanup", nil) })
}

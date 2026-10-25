package integration

import (
	"testing"

	"github.com/ozontech/testo"
)

type invalidContext struct {
	*testo.T
	Count int
}

type recursivePart struct {
	Next *recursivePart
}

type recursiveContext struct {
	*testo.T
	Part *recursivePart
}

func TestConstruction(t *testing.T) {
	testo.RunTest(t, func(*invalidContext) {})
	testo.RunTest(t, func(*recursiveContext) {})
}

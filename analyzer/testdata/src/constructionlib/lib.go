package constructionlib

import "github.com/ozontech/testo"

type Invalid struct{ Bad int }
type Cycle struct{ Next *Cycle }
type Context struct {
	*testo.T
	Bad string
}

package message_links

import "github.com/ozontech/testo"

type Bad struct{ testo.Suite[*testo.T] }

func (Bad) TestBad() {} // want "TESTO001: .* \\(https://github.com/nashabanov/testolint/blob/main/docs/rules/TESTO001.md\\)$"
type Good struct{ testo.Suite[*testo.T] }

func (Good) TestGood(t *testo.T) {}

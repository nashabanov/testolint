package integration

import (
	"github.com/ozontech/testo"
	"os"
)

type PrematureFile struct {
	testo.Suite[*testo.T]
	file *os.File
}

func (s *PrematureFile) TestFile(t *testo.T) {
	testo.Run(t, "child", func(child *testo.T) { child.Parallel(); _, _ = s.file.Read(nil) })
}
func (s *PrematureFile) AfterEach(*testo.T) { _ = s.file.Close() }

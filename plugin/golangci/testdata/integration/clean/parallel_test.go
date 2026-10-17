package clean

import (
	"github.com/ozontech/testo"
	"testing"
)

type Serial struct {
	testo.Suite[*testo.T]
	value int
	ctx   *testo.T
}

func (s *Serial) BeforeEach(t *testo.T) { s.value = 1; s.ctx = t; t.Cleanup(func() { t.Log(s.value) }) }
func (s *Serial) TestA(t *testo.T)      { t.Log(s.value); s.ctx.Log("ok") }
func (s *Serial) TestB(t *testo.T)      { t.Log(s.value) }

type Captured struct {
	testo.Suite[*testo.T]
	value int
}

func (s *Captured) BeforeEach(t *testo.T) { s.value = 1; v := s.value; t.Cleanup(func() { t.Log(v) }) }
func (s *Captured) TestA(t *testo.T)      { v := s.value; t.Parallel(); t.Log(v) }
func (s *Captured) TestB(t *testo.T)      { t.Parallel() }

func TestSafeParallelPatterns(t *testing.T) {
	t.Run("serial", func(t *testing.T) { testo.RunSuite(t, &Serial{}) })
	t.Run("captured", func(t *testing.T) { testo.RunSuite(t, &Captured{}) })
}

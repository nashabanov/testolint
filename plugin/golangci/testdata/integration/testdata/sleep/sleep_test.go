package sleep

import (
	"github.com/ozontech/testo"
	"time"
)

type Suite struct{ testo.Suite[*testo.T] }

func (Suite) TestSleep(t *testo.T) { time.Sleep(0) }
func (Suite) TestValid(t *testo.T) {}

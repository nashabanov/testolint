package sleep

import (
	"github.com/ozontech/testo"
	"testing"
	"time"
)

type Suite struct{ testo.Suite[*testo.T] }

func (Suite) TestSleep(t *testo.T) { time.Sleep(0) }
func (Suite) TestValid(t *testo.T) {}

func (Suite) Testfoo(t *testo.T) { time.Sleep(0) }

func TestMalformedName(t *testing.T) { testo.RunSuite(t, &Suite{}) }

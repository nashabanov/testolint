package no_sleep

import (
	"github.com/ozontech/testo"
	"time"
)

type S struct{ testo.Suite[*testo.T] }

func (*S) TestSleep(t *testo.T)  { time.Sleep(0) }       // want "TESTO027"
func (*S) BeforeEach(t *testo.T) { defer time.Sleep(0) } // want "TESTO027"
func helper()                    { time.Sleep(0) }
func (*S) TestGood(t *testo.T)   { helper() }

func callbacks(t *testo.T) {
	testo.Run(t, "child", func(t *testo.T) { time.Sleep(0) }) // want "TESTO027"
	testo.RunTest(t, func(t *testo.T) { time.Sleep(0) })      // want "TESTO027"
	_ = testo.Test(func(t *testo.T) { time.Sleep(0) })        // want "TESTO027"
}
func (*S) AfterEach(t *testo.T) { func() { time.Sleep(0) }() } // want "TESTO027"
func (*S) BeforeAll(t *testo.T) { go time.Sleep(0) }           // want "TESTO027"
func (*S) AfterAll(t *testo.T)  { time.Sleep(0) }              // want "TESTO027"
type Other struct{}

func (Other) TestOther(t *testo.T) { time.Sleep(0) }

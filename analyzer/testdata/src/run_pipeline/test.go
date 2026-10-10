package run_pipeline

import "github.com/ozontech/testo"

type T = *testo.T
type S struct{ testo.Suite[T] }

func (S) TestBad()          {}                  // want "TESTO001"
func (S) BeforeEach()       {}                  // want "TESTO002"
func (S) Testlower(t T)     {}                  // want "TESTO007"
func (S) Caseslower() []int { return []int{1} } // want "TESTO008"
func (S) TestParams(t T, p struct {
	Missing               int // want "TESTO003"
	private               int // want "TESTO009"
	Wrong, Invalid, Empty int
}) {
}
func (S) CasesWrong() []string     { return []string{"x"} } // want "TESTO004"
func (S) CasesInvalid(x int) []int { return []int{1} }      // want "TESTO005"
func (S) CasesUnused() []int       { return []int{1} }      // want "TESTO006"
func (S) CasesEmpty() []int        { return nil }           // want "TESTO010"

type Empty struct{ testo.Suite[T] } // want `TESTO011: suite "Empty" contains no tests`

func runs(t T) {
	// S's malformed hook makes S invalid as a run argument. Use Empty and the
	// base suite instead; suite checks still run for S in the same pass.
	testo.RunSuite(t, Empty{})             // want "TESTO_RUN_PIPELINE: run" "TESTO_RUN_PIPELINE: other"
	testo.RunSubSuite(t, testo.Suite[T]{}) // want "TESTO_RUN_PIPELINE: run" "TESTO_RUN_PIPELINE: other"
	RunSuite(t, Empty{})
}
func RunSuite(t T, s Empty) {}

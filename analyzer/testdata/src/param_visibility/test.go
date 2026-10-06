package param_visibility

import "github.com/ozontech/testo"

type T struct{ *testo.T }
type Suite struct{ testo.Suite[T] }

func (Suite) TestSimple(t T) {}

func (Suite) TestValid(t T, p struct {
	UserID int
	Role   string
}) {
}

func (Suite) CasesUserID() []int  { return []int{0} }
func (Suite) CasesRole() []string { return []string{"admin"} }

func (Suite) TestUnexported(t T, p struct {
	UserID int
	role   string // want `TESTO009: parameter field "role" must be exported`
}) {
}

func (Suite) TestMultiple(t T, p struct {
	userID int    // want `TESTO009: parameter field "userID" must be exported`
	role   string // want `TESTO009: parameter field "role" must be exported`
}) {
}

// An existing provider with an incompatible element must not cause TESTO004.
// Its malformed name still independently receives TESTO008.
func (Suite) CasesuserID() []string { return []string{"admin"} } // want "TESTO008"

type UnicodeParams struct {
	État string
	état string // want `TESTO009: parameter field "état" must be exported`
	中    int    // want `TESTO009: parameter field "中" must be exported`
}

func (Suite) TestUnicode(t T, p UnicodeParams) {}
func (Suite) CasesÉtat() []string              { return []string{"admin"} }

package operation_ignore

import a "github.com/ozontech/testo-allure"

type Client struct{}

func (Client) Read(...int)       {}
func receiverAction() Client     { return Client{} }
func assertionHelper(t T) Client { verify(); return Client{} }
func (*Suite) TestReceiver(t T) {
	receiverAction().Read() // want "TESTO029"
}
func (*Suite) TestReceiverPhase(t T) {
	verify()                // want "TESTO031"
	receiverAction().Read() // want "TESTO029" "TESTO030"
}
func (*Suite) TestReceiverHelper(t T) {
	assertionHelper(t).Read() // want "TESTO031"
}
func (*Suite) TestCoveredReceiver(t T) {
	a.Step(t, "work", func(t T) { receiverAction().Read(); assertionHelper(t).Read() })
}
func (*Suite) TestReceiverBeforeArguments(t T) {
	receiverAction().Read(assertionArgument()) // want "TESTO029" "TESTO031"
}
func assertionArgument() int { verify(); return 1 }
func (*Suite) TestUnusedReceiverClosure(t T) {
	_ = func() { receiverAction().Read() }
}

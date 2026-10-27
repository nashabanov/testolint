package providerbeforeeach

import "github.com/ozontech/testo"

type Suite struct {
	testo.Suite[*testo.T]
	value *int
}

func (s *Suite) BeforeEach(*testo.T)                       { s.value = new(int) }
func (*Suite) TestValue(t *testo.T, p struct{ Value int }) {}
func (s *Suite) CasesValue() []int {
	return []int{*s.value} // want "TESTO023:.*value.*BeforeEach"
}

type Alias = Suite

func run() {
	testo.RunSuite(nil, &Suite{})
	testo.RunSuite(nil, &Alias{}) // same provider: one diagnostic
	testo.RunSubSuite(&testo.T{}, &Suite{})
	value := 42
	testo.RunSuite(nil, &Suite{value: &value}) // constructor initialization is valid
	s := &Suite{}
	s.value = &value
	testo.RunSuite(nil, s) // local history is deliberately unknown
}

// BeforeAll runs before provider evaluation; dynamic providers are valid.
type BeforeAll struct {
	testo.Suite[*testo.T]
	value *int
}

func (s *BeforeAll) BeforeAll(*testo.T)                        { s.value = new(int) }
func (s *BeforeAll) BeforeEach(*testo.T)                       { s.value = new(int) }
func (*BeforeAll) TestValue(t *testo.T, p struct{ Value int }) {}
func (s *BeforeAll) CasesValue() []int                         { return []int{*s.value} }
func initialized()                                             { testo.RunSuite(nil, &BeforeAll{}) }

// A nil slice is valid data, not a required nonzero value.
type ZeroValue struct {
	testo.Suite[*testo.T]
	values []int
}

func (s *ZeroValue) BeforeEach(*testo.T)                       { s.values = []int{1} }
func (*ZeroValue) TestValue(t *testo.T, p struct{ Value int }) {}
func (s *ZeroValue) CasesValue() []int                         { return s.values }
func zeroValue()                                               { testo.RunSuite(nil, &ZeroValue{}) }

// A declaration alone does not prove that an instance starts with nil fields.
type Unused struct {
	testo.Suite[*testo.T]
	value *int
}

func (s *Unused) BeforeEach(*testo.T)                       { s.value = new(int) }
func (*Unused) TestValue(t *testo.T, p struct{ Value int }) {}
func (s *Unused) CasesValue() []int                         { return []int{*s.value} }

type IntPointer = *int
type Context = *testo.T
type Aliased struct {
	testo.Suite[Context]
	value IntPointer
}

func (s *Aliased) BeforeEach(Context)                       { s.value = new(int) }
func (*Aliased) TestValue(t Context, p struct{ Value int }) {}
func (s *Aliased) CasesValue() []int {
	return []int{*s.value} // want "TESTO023:"
}
func aliased() { testo.RunSuite(nil, &Aliased{}) }

// Value hooks write their own copy; promoted methods are unknown.
type ValueReceiver struct {
	testo.Suite[*testo.T]
	value *int
}

func (s ValueReceiver) BeforeEach(*testo.T)                       { s.value = new(int) }
func (ValueReceiver) TestValue(t *testo.T, p struct{ Value int }) {}
func (s *ValueReceiver) CasesValue() []int                        { return []int{*s.value} }
func valueReceiver()                                              { testo.RunSuite(nil, &ValueReceiver{}) }

type Promoted struct{ *Unused }

func promoted() { testo.RunSuite(nil, &Promoted{}) }

type Invalid struct {
	testo.Suite[*testo.T]
	value *int
}

func (s *Invalid) BeforeEach(*testo.T)                           { s.value = new(int) }
func (*Invalid) TestValue(t *testo.T, p struct{ Value int }) int { return 0 } // want "TESTO001:"
func (s *Invalid) CasesValue() []int                             { return []int{*s.value} }
func invalid()                                                   { testo.RunSuite(nil, &Invalid{}) }

type OtherSuite struct{ value *int }

func (s *OtherSuite) BeforeEach(*testo.T)                       { s.value = new(int) }
func (*OtherSuite) TestValue(t *testo.T, p struct{ Value int }) {}
func (s *OtherSuite) CasesValue() []int                         { return []int{*s.value} }

// A nil embedded marker pointer fails before collection; avoid a dependent
// provider diagnostic on that unsupported initialization path.
type PointerMarker struct {
	*testo.Suite[*testo.T]
	value *int
}

func (s *PointerMarker) BeforeEach(*testo.T)                       { s.value = new(int) }
func (*PointerMarker) TestValue(t *testo.T, p struct{ Value int }) {}
func (s *PointerMarker) CasesValue() []int                         { return []int{*s.value} }
func pointerMarker()                                               { testo.RunSuite(nil, &PointerMarker{}) }

// Earlier parameter collection can initialize state too, even with default
// BeforeAll. Unknown sibling providers must prevent the nil-state proof.
type ProviderSetup struct {
	testo.Suite[*testo.T]
	value *int
}

func (s *ProviderSetup) BeforeEach(*testo.T)                   { s.value = new(int) }
func (*ProviderSetup) TestA(t *testo.T, p struct{ Init int })  {}
func (*ProviderSetup) TestB(t *testo.T, p struct{ Value int }) {}
func (s *ProviderSetup) CasesInit() []int                      { s.value = new(int); return []int{1} }
func (s *ProviderSetup) CasesValue() []int                     { return []int{*s.value} }
func providerSetup()                                           { testo.RunSuite(nil, &ProviderSetup{}) }

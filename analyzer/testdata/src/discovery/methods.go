package discovery

// Methods live in a different file and may use aliases as receivers.
func (Alias) TestValid(t T, p struct{ Age int }) {}
func (*Active) CasesAge() []int                  { return []int{0} }
func (Active) BeforeEach(t T)                    {}
func (*Active) TestInvalid()                     {}                  // want "TESTO001"
func (Active) AfterEach()                        {}                  // want "TESTO002"
func (Active) CasesUnused() []int                { return []int{0} } // want "TESTO006"

func (Generic[P]) TestInvalid() {}
func (Indirect) TestInvalid()   {} // want "TESTO001"
func (Unrelated) TestInvalid()  {}
func (NonStruct) TestInvalid()  {}

package discovery

// Methods live in a different file and may use aliases as receivers.
func (Alias) TestValid(t T, p struct{ Age int }) {}
func (*Active) CasesAge() []int                  { return nil }
func (Active) BeforeEach(t T)                    {}
func (*Active) TestInvalid()                     {}             // want "TESTO001"
func (Active) AfterEach()                        {}             // want "TESTO002"
func (Active) CasesUnused() []int                { return nil } // want "TESTO006"

func (Generic[P]) TestInvalid() {}
func (Indirect) TestInvalid()   {}
func (Unrelated) TestInvalid()  {}
func (NonStruct) TestInvalid()  {}

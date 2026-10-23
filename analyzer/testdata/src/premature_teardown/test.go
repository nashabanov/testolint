package premature_teardown

import (
	tt "github.com/ozontech/testo"
	"os"
)

type File = os.File
type Context = tt.T

type Broken struct {
	tt.Suite[*Context]
	file  *File
	other *File
}

func (s *Broken) TestFile(t *Context) {
	tt.Run(t, "child", func(child *Context) { child.Parallel(); _, _ = s.file.Read(nil) })
}
func (s *Broken) AfterEach(t *Context) { _ = s.file.Close() } // want "TESTO020"

type Value struct {
	tt.Suite[*Context]
	file  *File
	other *File
}

func (s Value) TestFile(t *Context) {
	tt.Run(t, "child", func(child *Context) { child.Parallel(); _, _ = s.file.Write(nil) })
}
func (s Value) AfterEach(t *Context) { s.file.Close() } // want "TESTO020"

type Stat struct {
	tt.Suite[*Context]
	file  *File
	other *File
}

func (s *Stat) TestFile(t *Context) {
	tt.Run(t, "child", func(child *Context) { child.Parallel(); _, _ = s.file.Stat() })
}
func (s *Stat) AfterEach(t *Context) { _ = s.file.Close() } // want "TESTO020"

type Serial struct {
	tt.Suite[*Context]
	file  *File
	other *File
}

func (s *Serial) TestFile(t *Context) {
	tt.Run(t, "child", func(child *Context) { _, _ = s.file.Read(nil) })
}
func (s *Serial) AfterEach(t *Context) { _ = s.file.Close() }

type Other struct {
	tt.Suite[*Context]
	file  *File
	other *File
}

func (s *Other) TestFile(t *Context) {
	tt.Run(t, "child", func(child *Context) { child.Parallel(); _, _ = s.file.Read(nil) })
}
func (s *Other) AfterEach(t *Context) { _ = s.other.Close() }

type Cleanup struct {
	tt.Suite[*Context]
	file  *File
	other *File
}

func (s *Cleanup) TestFile(t *Context) {
	tt.Run(t, "child", func(child *Context) { child.Parallel(); _, _ = s.file.Read(nil) })
}
func (s *Cleanup) AfterEach(t *Context) { t.Cleanup(func() { _ = s.file.Close() }) }

type Replacement struct {
	tt.Suite[*Context]
	file  *File
	other *File
}

func (s *Replacement) TestFile(t *Context) {
	tt.Run(t, "child", func(child *Context) { child.Parallel(); _, _ = s.file.Read(nil) })
}
func (s *Replacement) AfterEach(t *Context) { _ = s.file.Close(); s.file = s.other }

type Before struct {
	tt.Suite[*Context]
	file  *File
	other *File
}

func (s *Before) TestFile(t *Context) {
	tt.Run(t, "child", func(child *Context) { child.Parallel(); _, _ = s.file.Read(nil) })
}
func (s *Before) AfterEach(t *Context) { _ = s.file.Close() }

func (s *Before) BeforeEach(*Context) { s.file = s.other }

type TwoTests struct {
	tt.Suite[*Context]
	file  *File
	other *File
}

func (s *TwoTests) TestFile(t *Context) {
	tt.Run(t, "child", func(child *Context) { child.Parallel(); _, _ = s.file.Read(nil) })
}
func (s *TwoTests) AfterEach(t *Context) { _ = s.file.Close() }
func (s *TwoTests) TestOther(*Context)   {}

type WrongContext struct {
	tt.Suite[*Context]
	file  *File
	other *File
}

func (s *WrongContext) TestFile(t *Context) {
	tt.Run(t, "child", func(child *Context) { t.Parallel(); _, _ = s.file.Read(nil) })
}
func (s *WrongContext) AfterEach(t *Context) { _ = s.file.Close() }

type Nested struct {
	tt.Suite[*Context]
	file  *File
	other *File
}

func (s *Nested) TestFile(t *Context) {
	tt.Run(t, "join", func(group *Context) {
		tt.Run(group, "child", func(child *Context) { child.Parallel(); _, _ = s.file.Read(nil) })
	})
}
func (s *Nested) AfterEach(t *Context) { _ = s.file.Close() }

type Aliased struct {
	tt.Suite[*Context]
	file  *File
	other *File
}

func (s *Aliased) TestFile(t *Context) {
	tt.Run(t, "child", func(child *Context) { child.Parallel(); f := s.file; _, _ = f.Read(nil) })
}
func (s *Aliased) AfterEach(t *Context) { _ = s.file.Close() }

type Branch struct {
	tt.Suite[*Context]
	file  *File
	other *File
}

func (s *Branch) TestFile(t *Context) {
	tt.Run(t, "child", func(child *Context) {
		child.Parallel()
		if s.file != nil {
			_, _ = s.file.Read(nil)
		}
	})
}
func (s *Branch) AfterEach(t *Context) { _ = s.file.Close() }

type BeforeParallel struct {
	tt.Suite[*Context]
	file  *File
	other *File
}

func (s *BeforeParallel) TestFile(t *Context) {
	tt.Run(t, "child", func(child *Context) { _, _ = s.file.Read(nil); child.Parallel() })
}
func (s *BeforeParallel) AfterEach(t *Context) { _ = s.file.Close() }

type fakeFile struct{}

func (*fakeFile) Close() error             { return nil }
func (*fakeFile) Read([]byte) (int, error) { return 0, nil }

type Unrelated struct {
	tt.Suite[*Context]
	file  *fakeFile
	other *File
}

func (s *Unrelated) TestFile(t *Context) {
	tt.Run(t, "child", func(child *Context) { child.Parallel(); _, _ = s.file.Read(nil) })
}
func (s *Unrelated) AfterEach(t *Context) { _ = s.file.Close() }

// Promoted method bodies and nested field ownership are outside this model.
type Promoted struct{ Broken }
type wrapper struct{ *File }
type Wrapped struct {
	tt.Suite[*Context]
	file *wrapper
}

func (s *Wrapped) TestFile(t *Context) {
	tt.Run(t, "child", func(child *Context) { child.Parallel(); _, _ = s.file.Read(nil) })
}
func (s *Wrapped) AfterEach(*Context) { _ = s.file.Close() }

type InvalidHook struct {
	tt.Suite[*Context]
	file *File
}

func (s *InvalidHook) TestFile(t *Context) {
	tt.Run(t, "child", func(child *Context) { child.Parallel(); _, _ = s.file.Read(nil) })
}
func (s *InvalidHook) AfterEach() { _ = s.file.Close() } // want "TESTO002"

type NotSuite struct{ file *File }

func (s *NotSuite) TestFile(t *Context) {
	tt.Run(t, "child", func(child *Context) { child.Parallel(); _, _ = s.file.Read(nil) })
}
func (s *NotSuite) AfterEach(*Context) { _ = s.file.Close() }

type UnknownArgument struct {
	tt.Suite[*Context]
	file *File
}

func (s *UnknownArgument) buffer() []byte { s.file = nil; return nil }
func (s *UnknownArgument) TestFile(t *Context) {
	tt.Run(t, "child", func(child *Context) { child.Parallel(); _, _ = s.file.Read(s.buffer()) })
}
func (s *UnknownArgument) AfterEach(*Context) { _ = s.file.Close() }

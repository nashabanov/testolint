package executionruntime

import (
	"errors"
	"github.com/ozontech/testo"
	"os"
	"testing"
)

func TestNilRunTest(t *testing.T) { testo.RunTest[*testo.T](t, nil) }
func TestNilAdapter(t *testing.T) { testo.Test[*testo.T](nil)(t) }
func TestCleanupRun(t *testing.T) {
	testo.RunTest(t, func(ctx *testo.T) {
		ctx.Cleanup(func() { testo.Run(ctx, "bad", nil) })
	})
}

type premature struct {
	testo.Suite[*testo.T]
	closed bool
}

func (s *premature) TestResource(t *testo.T) {
	testo.Run(t, "child", func(t *testo.T) {
		t.Parallel()
		if !s.closed {
			t.Fatal("AfterEach must precede parallel child resumption")
		}
		t.Log("AfterEach ran before parallel child resumed")
	})
}
func (s *premature) AfterEach(*testo.T) { s.closed = true }
func TestPremature(t *testing.T)        { testo.RunSuite(t, &premature{}) }

// Instrument each operation to verify the concrete os.File contract after the
// lifecycle teardown. Assertions intentionally exceed the analyzer's subset.
type closedFile struct {
	testo.Suite[*testo.T]
	file *os.File
}

func (s *closedFile) TestFile(t *testo.T) {
	testo.Run(t, "child", func(child *testo.T) {
		child.Parallel()
		_, readErr := s.file.Read(nil)
		_, writeErr := s.file.Write(nil)
		_, statErr := s.file.Stat()
		if !errors.Is(readErr, os.ErrClosed) || !errors.Is(writeErr, os.ErrClosed) || !errors.Is(statErr, os.ErrClosed) {
			child.Fatal("expected all three operations to observe os.ErrClosed")
		}
		child.Fatal("TESTO020 confirmed: Read, Write and Stat used the file after AfterEach closed it")
	})
}
func (s *closedFile) AfterEach(*testo.T) { _ = s.file.Close() }
func TestClosedFile(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "resource")
	if err != nil {
		t.Fatal(err)
	}
	testo.RunSuite(t, &closedFile{file: f})
}

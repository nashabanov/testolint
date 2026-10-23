package clean

import (
	"github.com/ozontech/testo"
	"os"
	"testing"
)

type CleanupFile struct {
	testo.Suite[*testo.T]
	file *os.File
}

func (s *CleanupFile) TestFile(t *testo.T) {
	testo.Run(t, "child", func(child *testo.T) {
		child.Parallel()
		if _, err := s.file.Stat(); err != nil {
			child.Fatal(err)
		}
	})
}
func (s *CleanupFile) AfterEach(t *testo.T) { t.Cleanup(func() { _ = s.file.Close() }) }

type SerialFile struct {
	testo.Suite[*testo.T]
	file *os.File
}

func (s *SerialFile) TestFile(t *testo.T) {
	testo.Run(t, "child", func(child *testo.T) {
		if _, err := s.file.Stat(); err != nil {
			child.Fatal(err)
		}
	})
}
func (s *SerialFile) AfterEach(t *testo.T) { _ = s.file.Close() }
func TestSafeFileTeardown(t *testing.T) {
	for _, parallel := range []bool{false, true} {
		t.Run(map[bool]string{false: "serial", true: "cleanup"}[parallel], func(t *testing.T) {
			f, err := os.CreateTemp(t.TempDir(), "resource")
			if err != nil {
				t.Fatal(err)
			}
			if parallel {
				testo.RunSuite(t, &CleanupFile{file: f})
			} else {
				testo.RunSuite(t, &SerialFile{file: f})
			}
		})
	}
}

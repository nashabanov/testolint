package clean

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/ozontech/testo"
	"github.com/ozontech/testo/testoplugin"
)

type parentEvidence struct{ root bool }

func (p *parentEvidence) Plugin(parent testoplugin.Plugin, _ ...testoplugin.Option) testoplugin.Spec {
	requireParentInterface(parent)
	// The interface contains a concrete type even at the suite root.
	previous := parent.(*parentEvidence)
	p.root = previous == nil
	return testoplugin.Spec{}
}

func TestProviderRuntime(t *testing.T) {
	cmd := exec.Command("go", "test", "../testdata/providerruntime", "-run", "^TestLateProvider$", "-count=1")
	out, err := cmd.CombinedOutput()
	if err == nil || !strings.Contains(string(out), "invalid memory address or nil pointer dereference") || !strings.Contains(string(out), "CasesValue") {
		t.Fatalf("expected provider nil dereference: %v\n%s", err, out)
	}
	cmd = exec.Command("go", "test", "../testdata/providerruntime", "-run", "^TestConstructorProvider$", "-count=1")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("initialized provider: %v\n%s", err, out)
	}
}

type parentEvidenceT struct {
	*testo.T
	Evidence *parentEvidence
}

func TestPluginParent(t *testing.T) {
	// Observe both the suite root and per-test parent supplied by Testo.
	testo.RunSuite(t, parentEvidenceSuite{})
}

type parentEvidenceSuite struct{ testo.Suite[*parentEvidenceT] }

func (parentEvidenceSuite) BeforeAll(t *parentEvidenceT) {
	if !t.Evidence.root {
		t.Fatal("suite parent must contain a typed nil")
	}
}
func (parentEvidenceSuite) TestChild(t *parentEvidenceT) {
	if t.Evidence.root {
		t.Fatal("test parent must be the suite plugin")
	}
}

type dynamicCasesSuite struct {
	testo.Suite[*testo.T]
	values []int
}

func (s *dynamicCasesSuite) BeforeAll(*testo.T)  { s.values = []int{42} }
func (s *dynamicCasesSuite) BeforeEach(*testo.T) { s.values = []int{99} }
func (s *dynamicCasesSuite) CasesValue() []int   { return s.values }
func (*dynamicCasesSuite) TestValue(t *testo.T, p struct{ Value int }) {
	if p.Value != 42 {
		t.Fatal("providers must observe BeforeAll, not BeforeEach")
	}
}

func TestDynamicProviders(t *testing.T) {
	testo.RunSuite(t, &dynamicCasesSuite{})
}

func requireParentInterface(parent testoplugin.Plugin) {
	if parent == nil {
		panic("Testo supplied an untyped nil parent")
	}
}

func TestEmptyMakeLength(t *testing.T) {
	values := make([]int, 0, 4)
	if len(values) != 0 || cap(values) != 4 {
		t.Fatal("invalid empty make evidence")
	}
}

// Run in a separate process: Options intentionally persists for future runs.
func TestIntentionalOptions(t *testing.T) {
	cmd := exec.Command("go", "test", "../testdata/optionsruntime", "-run", "^TestIntentionalOptions$", "-count=1")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("intentional future configuration: %v\n%s", err, out)
	}
}

func TestEmptyMakeProvider(t *testing.T) {
	cmd := exec.Command("go", "test", "../testdata/providerruntime", "-run", "^TestEmptyMakeProvider$", "-count=1", "-v")
	out, err := cmd.CombinedOutput()
	if err != nil || !strings.Contains(string(out), "returned an empty slice") || !strings.Contains(string(out), "will not run") {
		t.Fatalf("empty provider warning/omission: %v\n%s", err, out)
	}
}

// Collection of an earlier parameterized method can establish initialization
// before the next provider, without using BeforeAll.
type providerSetupSuite struct {
	testo.Suite[*testo.T]
	value *int
}

func (s *providerSetupSuite) BeforeEach(*testo.T)                  { s.value = new(int) }
func (*providerSetupSuite) TestA(t *testo.T, p struct{ Init int }) {}
func (*providerSetupSuite) TestB(t *testo.T, p struct{ Value int }) {
	if p.Value != 42 {
		t.Fatal("earlier provider initialization was lost")
	}
}
func (s *providerSetupSuite) CasesInit() []int {
	value := 42
	s.value = &value
	return []int{1}
}
func (s *providerSetupSuite) CasesValue() []int { return []int{*s.value} }
func TestProviderInitialization(t *testing.T)   { testo.RunSuite(t, &providerSetupSuite{}) }

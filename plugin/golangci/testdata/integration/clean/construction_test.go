package clean

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/ozontech/testo"
	"github.com/ozontech/testo/testoplugin"
)

type constructionLeaf struct{ Value *int }
type constructionPointer *constructionLeaf
type constructionCycle struct{ Next *constructionCycle }
type constructionPlugin struct{ Value *int }

func (*constructionPlugin) Plugin(testoplugin.Plugin, ...testoplugin.Option) testoplugin.Spec {
	return testoplugin.Spec{}
}

type constructionT struct {
	*testo.T
	Plugin *constructionPlugin
	First  *constructionLeaf
	Second *constructionLeaf
	Deep   **constructionLeaf
	Named  constructionPointer
	Slice  *[]*constructionCycle
	Map    *map[string]*constructionCycle
	Array  *[1]*constructionCycle
	Any    *any
	hidden *constructionCycle
	value  int
}

type constructionAlias = constructionT
type constructionSuite struct {
	testo.Suite[*constructionAlias]
}

func (constructionSuite) TestFields(ctx *constructionAlias) {
	if ctx.Plugin == nil || ctx.Plugin.Value == nil || ctx.First == nil || ctx.Second == nil ||
		ctx.First == ctx.Second || ctx.Deep == nil || *ctx.Deep == nil || ctx.Named == nil || ctx.Slice == nil ||
		ctx.Map == nil || ctx.Array == nil || ctx.Any == nil || ctx.hidden != nil || ctx.value != 0 {
		ctx.Fatal("unexpected construction graph")
	}
}

func TestConstructionSafe(t *testing.T) {
	testo.RunSuite(t, constructionSuite{})
	testo.RunSuite(t, &constructionSuite{})
	testo.RunTest(t, func(ctx constructionT) {
		if ctx.Plugin == nil || ctx.First == nil {
			ctx.Fatal("value context was not initialized")
		}
	})
	t.Run("adapter", testo.Test(func(ctx *constructionT) {
		testo.Run(ctx, "child", func(child *constructionT) {
			if child.Plugin == nil || child.First == nil {
				child.Fatal("nested context was not initialized")
			}
		})
	}))
}

func TestConstructionRuntime(t *testing.T) {
	for _, tc := range []struct{ name, evidence string }{
		{"TestInvalidField", "all exported fields in T and plugins must be pointers: field constructionruntime.invalidContext.Count"},
		{"TestInvalidPluginField", "all exported fields in T and plugins must be pointers: field constructionruntime.invalidPlugin.Value"},
		{"TestRecursive", "infinite type recursion detected for constructionruntime.recursive"},
		{"TestMutual", "infinite type recursion detected for constructionruntime.left"},
		{"TestValueRecursive", "infinite type recursion detected for constructionruntime.valueRecursive"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cmd := exec.Command("go", "test", "../testdata/constructionruntime", "-run", "^"+tc.name+"$", "-count=1", "-v")
			out, err := cmd.CombinedOutput()
			if err == nil || !strings.Contains(string(out), tc.evidence) {
				t.Fatalf("expected construction failure %q: %v\n%s", tc.evidence, err, out)
			}
		})
	}
}

package analyzer

import (
	"go/types"

	"golang.org/x/tools/go/analysis"
)

func checkCasesSignature(p *analysis.Pass, suite *Suite) {
	for _, cases := range suite.Cases {
		if _, problem := casesSignature(cases.Func.Type()); problem != "" {
			p.Reportf(cases.Pos(suite), "TESTO005: %s", problem)
		}
	}
}

// casesSignature returns the provider's slice or the reason its signature is
// invalid. A non-signature type has no slice and no diagnostic reason.
func casesSignature(t types.Type) (*types.Slice, string) {
	sig, ok := t.(*types.Signature)
	if !ok {
		return nil, ""
	}
	if sig.Params().Len() != 0 {
		return nil, "cases provider must not accept parameters"
	}
	if sig.Results().Len() != 1 {
		return nil, "cases provider must return exactly one slice"
	}
	slice, ok := sig.Results().At(0).Type().Underlying().(*types.Slice)
	if !ok {
		return nil, "cases provider must return a slice"
	}
	return slice, ""
}

func isSliceType(t types.Type) bool {
	_, ok := t.Underlying().(*types.Slice)
	return ok
}

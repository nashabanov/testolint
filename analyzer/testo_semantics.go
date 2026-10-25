package analyzer

import (
	"go/types"
	"strings"
	"unicode"
	"unicode/utf8"
)

const testoImportPath = "github.com/ozontech/testo"

func isTestoSuiteType(t types.Type) bool {
	t = types.Unalias(t)
	if ptr, ok := t.(*types.Pointer); ok {
		t = types.Unalias(ptr.Elem())
	}

	named, ok := t.(*types.Named)
	if !ok {
		return false
	}

	obj := named.Obj()
	if obj == nil || obj.Pkg() == nil {
		return false
	}

	return obj.Pkg().Path() == testoImportPath &&
		obj.Name() == "Suite"
}

// testoTType uses Testo's private suite-interface marker rather than
// reconstructing Go's embedding and promotion rules. The selected method's
// receiver is the instantiated testo.Suite[T], even for indirect embedding.
func testoTType(methodSet *types.MethodSet) types.Type {
	for selection := range methodSet.Methods() {
		f := selection.Obj().(*types.Func)
		if f.Name() != "private" || f.Pkg() == nil || f.Pkg().Path() != testoImportPath {
			continue
		}
		sig := f.Type().(*types.Signature)
		if !isTestoSuiteType(sig.Recv().Type()) {
			continue
		}
		t := types.Unalias(sig.Recv().Type())
		if ptr, ok := t.(*types.Pointer); ok {
			t = types.Unalias(ptr.Elem())
		}
		args := t.(*types.Named).TypeArgs()
		if args.Len() == 1 {
			return args.At(0)
		}
	}
	return nil
}

func standardTestoT(t types.Type) bool {
	ptr, ok := types.Unalias(t).(*types.Pointer)
	if !ok {
		return false
	}
	named, ok := types.Unalias(ptr.Elem()).(*types.Named)
	return ok && named.Obj().Pkg() != nil &&
		named.Obj().Pkg().Path() == testoImportPath && named.Obj().Name() == "T"
}

func isStructType(t types.Type) bool {
	_, ok := t.Underlying().(*types.Struct)
	return ok
}

func isSliceType(t types.Type) bool {
	_, ok := t.Underlying().(*types.Slice)
	return ok
}

func isHookMethod(name string) bool {
	switch name {
	case "BeforeAll", "BeforeEach", "AfterEach", "AfterAll":
		return true
	default:
		return false
	}
}

func casesName(method *method) string {
	return strings.TrimPrefix(
		method.fn.Name(),
		"Cases",
	)
}

func isValidPrefixedName(name, prefix string) bool {
	if !strings.HasPrefix(name, prefix) {
		return false
	}

	if len(name) == len(prefix) {
		return true
	}

	r, _ := utf8.DecodeRuneInString(name[len(prefix):])
	return !unicode.IsLower(r)
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

func validParallelSignature(sig *types.Signature, context types.Type, name string) bool {
	if sig.Variadic() || sig.Results().Len() != 0 {
		return false
	}
	switch sig.Params().Len() {
	case 1:
	case 2:
		if !isValidPrefixedName(name, "Test") || !isStructType(sig.Params().At(1).Type()) {
			return false
		}
	default:
		return false
	}
	return types.Identical(sig.Params().At(0).Type(), context)
}

// validParams reports whether parameter checks have a meaningful runtime target.
func (t *test) validParams(suite *suite) bool {
	sig := t.fn().Type().(*types.Signature)
	return isValidPrefixedName(t.fn().Name(), "Test") &&
		sig.Params().Len() == 2 && !sig.Variadic() && sig.Results().Len() == 0 &&
		types.Identical(sig.Params().At(0).Type(), suite.tType) &&
		isStructType(sig.Params().At(1).Type())
}

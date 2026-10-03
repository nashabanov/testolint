package analyzer

import "testing"

func TestIsValidPrefixedName(t *testing.T) {
	for _, prefix := range []string{"Test", "Cases"} {
		for _, tc := range []struct {
			suffix string
			want   bool
		}{
			{"", true},
			{"Foo", true},
			{"1", true},
			{"_Foo", true},
			{"foo", false},
			{"É", true},
			{"é", false},
			{"中", true},
			{"ǅ", true},
		} {
			name := prefix + tc.suffix
			t.Run(name, func(t *testing.T) {
				if got := isValidPrefixedName(name, prefix); got != tc.want {
					t.Errorf("isValidPrefixedName(%q, %q) = %v, want %v", name, prefix, got, tc.want)
				}
			})
		}
		for _, name := range []string{"", "Te", "OtherFoo", "testFoo", "casesFoo"} {
			t.Run(prefix+"/"+name, func(t *testing.T) {
				if isValidPrefixedName(name, prefix) {
					t.Errorf("isValidPrefixedName(%q, %q) = true, want false", name, prefix)
				}
			})
		}
	}
}

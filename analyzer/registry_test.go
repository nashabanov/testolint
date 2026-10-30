package analyzer

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

func TestRegistry(t *testing.T) {
	registered := make(map[ruleID]bool)
	names := make(map[string]bool)
	checkers := make(map[uintptr]bool)
	for _, registration := range registrations {
		count := 0
		for _, checker := range []any{registration.suite, registration.parallel, registration.afterParallel, registration.run, registration.packageCheck} {
			value := reflect.ValueOf(checker)
			if value.IsNil() {
				continue
			}
			count++
			pointer := value.Pointer()
			if checkers[pointer] {
				t.Error("checker registered twice")
			}
			checkers[pointer] = true
		}
		if count != 1 || len(registration.ids) == 0 {
			t.Errorf("invalid checker registration: %+v", registration)
		}
		for _, id := range registration.ids {
			if registered[id] {
				t.Errorf("duplicate registration: %s", id)
			}
			registered[id] = true
			if _, ok := registry[id]; !ok {
				t.Errorf("missing metadata: %s", id)
			}
		}
	}
	readme, err := os.ReadFile("../README.md")
	if err != nil {
		t.Fatal(err)
	}
	documented := make(map[ruleID]bool)
	for _, match := range regexp.MustCompile("(?m)^\\| `([A-Z]+[0-9]+)` ").FindAllSubmatch(readme, -1) {
		documented[ruleID(match[1])] = true
	}
	for id, meta := range registry {
		if !registered[id] {
			t.Errorf("unregistered rule: %s", id)
		}
		if !documented[id] {
			t.Errorf("undocumented rule: %s", id)
		}
		delete(documented, id)
		if meta.Name == "" || names[meta.Name] {
			t.Errorf("empty or duplicate name: %s", id)
		}
		names[meta.Name] = true
		if meta.Class != correctness && meta.Class != policy {
			t.Errorf("invalid class: %s", id)
		}
		if !meta.EnableByDefault {
			t.Errorf("existing rule disabled by default: %s", id)
		}
		if meta.DocURL != ruleDocURL {
			t.Errorf("invalid documentation URL: %s", id)
		}
	}
	for id := range documented {
		t.Errorf("documented rule missing metadata: %s", id)
	}
	// Keep every existing diagnostic covered without depending on map iteration order.
	for i := 1; i <= 24; i++ {
		if !registered[ruleID(fmt.Sprintf("TESTO%03d", i))] {
			t.Errorf("missing existing rule %d", i)
		}
	}
	files, err := filepath.Glob("rule_*.go")
	if err != nil {
		t.Fatal(err)
	}
	emitted := make(map[ruleID]bool)
	for _, file := range files {
		if strings.HasSuffix(file, "_test.go") {
			continue
		}
		source, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		for _, id := range regexp.MustCompile(`"TESTO[0-9]{3}"`).FindAll(source, -1) {
			emitted[ruleID(strings.Trim(string(id), "\""))] = true
		}
	}
	for id := range emitted {
		if !registered[id] {
			t.Errorf("emitted diagnostic missing registration: %s", id)
		}
	}
	for id := range registered {
		if !emitted[id] {
			t.Errorf("registered diagnostic has no implementation: %s", id)
		}
	}
}

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
	for _, match := range regexp.MustCompile("(?m)^\\| \\[`([A-Z]+[0-9]+)`\\]").FindAllSubmatch(readme, -1) {
		documented[ruleID(match[1])] = true
	}
	examples, err := os.ReadFile("../docs/rules.md")
	if err != nil {
		t.Fatal(err)
	}
	for id, meta := range registry {
		if !registered[id] {
			t.Errorf("unregistered rule: %s", id)
		}
		defaultState := "disabled"
		if meta.EnableByDefault {
			defaultState = "enabled"
		}
		row := fmt.Sprintf("| [`%s`](docs/rules/%s.md) | `%s` | %s | %s |", id, id, meta.Name, meta.Class, defaultState)
		if !strings.Contains(string(readme), row) {
			t.Errorf("documentation metadata differs from registry: %s", id)
		}
		page, err := os.ReadFile("../docs/rules/" + string(id) + ".md")
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(page), "Violation") || !strings.Contains(string(page), "Valid") {
			t.Errorf("missing examples: %s", id)
		}
		if !strings.Contains(string(examples), "rules/"+string(id)+".md") {
			t.Errorf("missing index link: %s", id)
		}
		if !documented[id] {
			t.Errorf("undocumented rule: %s", id)
		}
		delete(documented, id)
		if meta.Name == "" || names[meta.Name] {
			t.Errorf("empty or duplicate name: %s", id)
		}
		names[meta.Name] = true
		if meta.Class != correctness && meta.Class != policy && meta.Class != style {
			t.Errorf("invalid class: %s", id)
		}
		if meta.EnableByDefault != (id != "TESTO025" && id != "TESTO026") {
			t.Errorf("unexpected default enablement: %s", id)
		}
		if meta.DocURL != "https://github.com/nashabanov/testolint/blob/main/docs/rules/"+string(id)+".md" {
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

// Published navigation entry points must continue to reach the per-rule pages.
func TestDocumentationLinks(t *testing.T) {
	files, err := filepath.Glob("../docs/rules/*.md")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != len(registry) {
		t.Fatalf("rule pages=%d, registry=%d", len(files), len(registry))
	}
	files = append(files, "../README.md", "../docs/rules.md", "../docs/contract.md", "../docs/analysis-scope.md")
	for _, file := range files {
		source, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		prose := regexp.MustCompile("(?s)```.*?```|`[^`]*`").ReplaceAllString(string(source), "")
		for _, match := range regexp.MustCompile(`\[[^\]]*\]\(([^)]+)\)`).FindAllStringSubmatch(prose, -1) {
			link, _, _ := strings.Cut(match[1], "#")
			if link == "" || strings.Contains(link, "://") {
				continue
			}
			if _, err := os.Stat(filepath.Join(filepath.Dir(file), link)); err != nil {
				t.Errorf("%s: broken link %s: %v", file, link, err)
			}
		}
	}
}

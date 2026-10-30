package analyzer

import (
	"fmt"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// Config selects diagnostics independently of checker execution. The zero value
// uses registry defaults. DisableAll starts with no enabled rules; Enable and
// Disable then apply explicit selections. Overlapping or duplicate IDs are errors.
type Config struct {
	DisableAll bool     `json:"disable-all"`
	Enable     []string `json:"enable"`
	Disable    []string `json:"disable"`
}

// New creates an independent analyzer with validated configuration and CLI flags.
// Neither config slices nor mutable flags are shared with another instance or
// the default Analyzer. Flags must be configured before analysis starts.
func New(config Config) (*analysis.Analyzer, error) {
	config.Enable = append([]string(nil), config.Enable...)
	config.Disable = append([]string(nil), config.Disable...)
	if _, err := resolveConfig(config); err != nil {
		return nil, err
	}
	a := &analysis.Analyzer{Name: Analyzer.Name, Doc: Analyzer.Doc}
	a.Flags.BoolVar(&config.DisableAll, "disable-all", config.DisableAll, "start with all rules disabled; enable selected IDs with -enable")
	a.Flags.Var(&selectionFlag{config: &config, enable: true}, "enable", "enable comma-separated rule IDs (for example TESTO006,TESTO019)")
	a.Flags.Var(&selectionFlag{config: &config}, "disable", "disable comma-separated rule IDs (must not overlap -enable)")
	a.Run = func(p *analysis.Pass) (any, error) {
		selected, err := resolveConfig(config)
		if err != nil {
			return nil, err
		}
		return runSelected(p, selected)
	}
	return a, nil
}

func resolveConfig(config Config) (map[ruleID]bool, error) {
	enable, err := validateIDs("enable", config.Enable)
	if err != nil {
		return nil, err
	}
	disable, err := validateIDs("disable", config.Disable)
	if err != nil {
		return nil, err
	}
	// Walk the input list for stable error ordering rather than iterating a map.
	for _, id := range config.Enable {
		if disable[ruleID(id)] {
			return nil, fmt.Errorf("testolint: rule %s appears in both enable and disable; remove it from one list", id)
		}
	}
	selected := make(map[ruleID]bool, len(registry))
	for id, meta := range registry {
		selected[id] = !config.DisableAll && meta.EnableByDefault
		if enable[id] {
			selected[id] = true
		}
		if disable[id] {
			selected[id] = false
		}
	}
	return selected, nil
}

func validateIDs(setting string, ids []string) (map[ruleID]bool, error) {
	seen := make(map[ruleID]bool, len(ids))
	for _, value := range ids {
		id := ruleID(value)
		if value == "" {
			return nil, fmt.Errorf("testolint: empty rule ID in %s; use IDs such as TESTO006", setting)
		}
		if _, ok := registry[id]; !ok {
			return nil, fmt.Errorf("testolint: unknown rule %q in %s; use a documented TESTO### ID", value, setting)
		}
		if seen[id] {
			return nil, fmt.Errorf("testolint: duplicate rule %s in %s; list each ID once", id, setting)
		}
		seen[id] = true
	}
	return seen, nil
}

type selectionFlag struct {
	config *Config
	enable bool
	set    bool
}

func (f *selectionFlag) String() string {
	if f.config == nil {
		return ""
	}
	if f.enable {
		return strings.Join(f.config.Enable, ",")
	}
	return strings.Join(f.config.Disable, ",")
}

func (f *selectionFlag) Set(value string) error {
	if f.set {
		return fmt.Errorf("testolint: selection flags may only be set once; combine IDs in one comma-separated list")
	}
	var ids []string
	if value != "" {
		ids = strings.Split(value, ",")
		for i := range ids {
			ids[i] = strings.TrimSpace(ids[i])
		}
	}
	candidate := *f.config
	if f.enable {
		candidate.Enable = ids
	} else {
		candidate.Disable = ids
	}
	if _, err := resolveConfig(candidate); err != nil {
		return err
	}
	*f.config = candidate
	f.set = true
	return nil
}

func runSelected(p *analysis.Pass, selected map[ruleID]bool) (any, error) {
	pass := *p
	pass.Report = func(d analysis.Diagnostic) {
		if selected[ruleID(d.Category)] {
			p.Report(d)
		}
	}
	return runWithRules(&pass, selectedRules(selected))
}

func selectedRules(selected map[ruleID]bool) ruleSet {
	var rules ruleSet
	for _, registration := range registrations {
		enabled := false
		for _, id := range registration.ids {
			enabled = enabled || selected[id]
		}
		if !enabled {
			continue
		}
		if registration.suite != nil {
			rules.suites = append(rules.suites, registration.suite)
		}
		if registration.parallel != nil {
			rules.parallel = append(rules.parallel, registration.parallel)
		}
		if registration.afterParallel != nil {
			rules.afterParallel = append(rules.afterParallel, registration.afterParallel)
		}
		if registration.run != nil {
			rules.runs = append(rules.runs, registration.run)
		}
		if registration.packageCheck != nil {
			rules.packages = append(rules.packages, registration.packageCheck)
		}
	}
	return rules
}

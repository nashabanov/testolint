package analyzer

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type styleFlags struct {
	config   *Config
	strict   *bool
	styleSet bool
}

type styleJSONFlag struct{ state *styleFlags }

func (f *styleJSONFlag) String() string { return "" }
func (f *styleJSONFlag) Set(value string) error {
	if f.state.styleSet {
		return fmt.Errorf("testolint: -style-config may only be set once")
	}
	decoder := json.NewDecoder(strings.NewReader(value))
	decoder.DisallowUnknownFields()
	var style *StyleConfig
	if err := decoder.Decode(&style); err != nil {
		return fmt.Errorf("testolint: invalid style config: %w", err)
	}
	if style == nil {
		return fmt.Errorf("testolint: style config must be an object")
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		return fmt.Errorf("testolint: trailing data in style config")
	}
	if f.state.strict != nil {
		style.Phases.Strict = f.state.strict
	}
	candidate := *f.state.config
	candidate.Style = *style
	if _, err := resolveConfig(candidate); err != nil {
		return err
	}
	*f.state.config = candidate
	f.state.styleSet = true
	return nil
}

type phaseStrictFlag struct{ state *styleFlags }

func (f *phaseStrictFlag) String() string {
	if f.state == nil || f.state.config.Style.Phases.Strict == nil {
		return "true"
	}
	return strconv.FormatBool(*f.state.config.Style.Phases.Strict)
}
func (*phaseStrictFlag) IsBoolFlag() bool { return true }
func (f *phaseStrictFlag) Set(value string) error {
	if f.state.strict != nil {
		return fmt.Errorf("testolint: -phase-order-strict may only be set once")
	}
	strict, err := strconv.ParseBool(value)
	if err != nil {
		return err
	}
	f.state.strict = &strict
	f.state.config.Style.Phases.Strict = &strict
	return nil
}

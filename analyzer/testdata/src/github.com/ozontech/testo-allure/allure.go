package allure

import "github.com/ozontech/testo"

type PluginAllure struct{ *testo.T }

func (*PluginAllure) Title(string)                            {}
func (*PluginAllure) Tags(...string)                          {}
func (*PluginAllure) Description(string)                      {}
func Step[T any](t T, name string, f func(T), options ...any) {}

func (a *PluginAllure) unwrapAllure() *PluginAllure { return a }
func (*PluginAllure) Titlef(string, ...any)         {}
func (*PluginAllure) Owner(string)                  {}

func Setup[T any](t T, name string, f func(T), options ...any)    {}
func TearDown[T any](t T, name string, f func(T), options ...any) {}

type Assertions struct{}
type Requirements struct{}

func (*PluginAllure) Assert() Assertions       { return Assertions{} }
func (*PluginAllure) Require() Requirements    { return Requirements{} }
func (Assertions) Equal(any, any, ...any) bool { return true }
func (Requirements) Equal(any, any, ...any)    {}
func (a Assertions) Masked() Assertions        { return a }
func (r Requirements) Hidden() Requirements    { return r }

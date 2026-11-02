package allure

import "github.com/ozontech/testo"

type PluginAllure struct{ *testo.T }

func (*PluginAllure) Title(string)                            {}
func (*PluginAllure) Tags(...string)                          {}
func (*PluginAllure) Description(string)                      {}
func Step[T any](t T, name string, f func(T), options ...any) {}

package testoplugin

type Option struct{ Value any }
type Spec struct{}
type Plugin interface {
	Plugin(Plugin, ...Option) Spec
}

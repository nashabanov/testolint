package integration

import (
	"github.com/ozontech/testo"
	"github.com/ozontech/testo/testoplugin"
)

type parentPlugin struct{}

func (*parentPlugin) Plugin(parent testoplugin.Plugin, _ ...testoplugin.Option) testoplugin.Spec {
	if parent == nil {
		return testoplugin.Spec{}
	}
	return testoplugin.Spec{}
}

type parentContext struct {
	*testo.T
	Parent *parentPlugin
}

func pluginParentRoot() {
	testo.RunTest(nil, func(*parentContext) {})
}

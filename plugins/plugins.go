package plugins

import (
	"github.com/open-policy-agent/opa/runtime"
	"github.com/thomasdarimont/custom-opa/custom-opa-spicedb/plugins/auhtzed"
)

func Register() {
	runtime.RegisterPlugin(auhtzed.PluginName, auhtzed.Factory{})
}

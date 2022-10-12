package plugins

import (
	"github.com/open-policy-agent/opa/runtime"
	"github.com/thomasdarimont/custom-opa/custom-opa-spicedb/plugins/authzed"
)

func Register() {
	runtime.RegisterPlugin(authzed.PluginName, authzed.Factory{})
}

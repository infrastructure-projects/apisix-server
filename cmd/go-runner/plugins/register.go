package plugins

import (
	"github.com/apache/apisix-go-plugin-runner/cmd/go-runner/plugins/security"
	"github.com/apache/apisix-go-plugin-runner/pkg/log"
	"github.com/apache/apisix-go-plugin-runner/pkg/plugin"
)

func init() {
	// security filter
	if err := plugin.RegisterPlugin(security.New()); err != nil {
		log.Fatalf("failed to register plugin SecurityFilter: %s", err)
		return
	}
}

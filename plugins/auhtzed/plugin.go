package auhtzed

import (
	"context"
	"github.com/authzed/authzed-go/v1"
	"github.com/authzed/grpcutil"
	"github.com/open-policy-agent/opa/plugins"
	"github.com/open-policy-agent/opa/util"
	"google.golang.org/grpc"
	"sync"
)

const PluginName = "authzed"

type Config struct {
	Endpoint string `json:"endpoint"`
	Insecure bool   `json:"insecure"`
	Token    string `json:"token"`
}

type AuthzedPlugin struct {
	manager *plugins.Manager
	mtx     sync.Mutex
	config  Config
	client  *authzed.Client
}

var instance *AuthzedPlugin = nil

func GetAuthzedClient() *authzed.Client {

	if instance == nil {
		return nil
	}

	instance.mtx.Lock()
	defer instance.mtx.Unlock()

	return instance.client
}

func (a *AuthzedPlugin) Start(ctx context.Context) error {

	grpcSecurity := grpcutil.WithSystemCerts(grpcutil.VerifyCA)
	if a.config.Insecure {
		grpcSecurity = grpc.WithInsecure()
	}

	client, err := authzed.NewClient(
		a.config.Endpoint,
		// grpcutil.WithSystemCerts(grpcutil.VerifyCA),
		grpcSecurity,
		grpcutil.WithInsecureBearerToken(a.config.Token),
	)

	a.client = client

	// HACK to expose plugin instance to be able to access the authzed client from the custom authzed check_permission builtin
	instance = a

	return err

}

func (a *AuthzedPlugin) Stop(ctx context.Context) {
}

func (a *AuthzedPlugin) Reconfigure(ctx context.Context, config any) {

	a.mtx.Lock()
	defer a.mtx.Unlock()

	if a.config.Endpoint != config.(Config).Endpoint {
		a.Stop(ctx)
		if err := a.Start(ctx); err != nil {
			a.manager.UpdatePluginStatus(PluginName, &plugins.Status{State: plugins.StateErr})
		}
	}
	a.config = config.(Config)
}

type Factory struct{}

func (Factory) New(m *plugins.Manager, config any) plugins.Plugin {

	m.UpdatePluginStatus(PluginName, &plugins.Status{State: plugins.StateNotReady})

	return &AuthzedPlugin{
		manager: m,
		config:  config.(Config),
	}
}

func (Factory) Validate(_ *plugins.Manager, config []byte) (any, error) {
	parsedConfig := Config{}
	err := util.Unmarshal(config, &parsedConfig)
	return parsedConfig, err
}

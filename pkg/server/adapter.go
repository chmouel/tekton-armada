package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/chmouel/armadas/pkg/app"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"go.uber.org/zap"

	"knative.dev/eventing/pkg/adapter/v2"
	"knative.dev/pkg/logging"
)

type envConfig struct {
	adapter.EnvConfig
}

func NewEnvConfig() adapter.EnvConfigAccessor {
	return &envConfig{}
}

type controller struct {
	Logger *zap.SugaredLogger
	App    *app.App
}

func (c *controller) Start(ctx context.Context) error {
	router := c.getRouter()
	portAddr := fmt.Sprintf("%s:%d", c.App.ServerAddress, c.App.Port)
	c.Logger.Infof("Starting Controller on %s", portAddr)
	//nolint:gosec
	return http.ListenAndServe(portAddr, router)
}

func New() adapter.AdapterConstructor {
	return func(ctx context.Context, _ adapter.EnvConfigAccessor, _ cloudevents.Client) adapter.Adapter {
		logger := logging.FromContext(ctx)
		return &controller{
			Logger: logger,
			App:    app.ParseFlags(),
		}
	}
}

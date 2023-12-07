package main

import (
	"github.com/chmouel/armadas/pkg/server"
	evadapter "knative.dev/eventing/pkg/adapter/v2"
	"knative.dev/pkg/injection"
	"knative.dev/pkg/signals"
)

func main() {
	cfg := injection.ParseAndGetRESTConfigOrDie()

	ctx := signals.NewContext()
	ctx = injection.WithConfig(ctx, cfg)

	evadapter.MainWithContext(ctx, "server", server.NewEnvConfig, server.New())
}

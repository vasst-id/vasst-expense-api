package main

import (
	"log"

	// "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"github.com/vasst-id/vasst-expense-api/config"
	"github.com/vasst-id/vasst-expense-api/internal/subscriber"
)

func main() {

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Error loading config: %s", err.Error())
	}

	// tracer.Start(
	// 	tracer.WithEnv(cfg.Env),
	// 	tracer.WithService(subscriber.ServiceName),
	// 	tracer.WithServiceVersion(cfg.SvcVersion),
	// 	// tracer.WithAgentAddr(cfg.DatadogAgentHost+":8126"),
	// 	tracer.WithPropagator(tracer.NewPropagator(&tracer.PropagatorConfig{B3: true})),
	// 	tracer.WithLogStartup(false),
	// )
	// defer tracer.Stop()

	subscriber.Run(cfg)
}

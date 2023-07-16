package tracer

import (
	"fmt"

	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

// Init global tracer
func InitGlobal(service string, host string, port string) error {
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: fmt.Sprintf("%s:%s", host, port),
		},
	}

	if _, err := cfg.InitGlobalTracer(service); err != nil {
		return err
	}

	return nil
}

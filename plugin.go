package pcx_frontend_api_rr

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	rrErrors "github.com/roadrunner-server/errors"
	httpapi "github.com/roadrunner-server/http/v5/api"
	"go.uber.org/zap"
)

// PluginName contains default service name.
const (
	PluginName     = "pcx_frontend_api_rr"
	RootPluginName = "http"
)

type Configurer interface {
	UnmarshalKey(name string, out any) error
	Has(name string) bool
}

type Logger interface {
	NamedLogger(name string) *zap.Logger
}

type Plugin struct {
	log      *zap.Logger
	fiberApp *fiber.App
}

func (p *Plugin) Init(cfg httpapi.Configurer, log httpapi.Logger) error {
	const op = rrErrors.Op("pcx_frontend_api_rr_init")

	if !cfg.Has(RootPluginName) {
		return rrErrors.E(op, rrErrors.Disabled)
	}

	p.log = zap.L().Named(PluginName)

	p.fiberApp = fiber.New()
	p.fiberApp.Get("/stub", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNoContent)
	})

	return nil
}

func (p *Plugin) Name() string {
	return PluginName
}

func (p *Plugin) Middleware(next http.Handler) http.Handler {
	p.fiberApp.Use(func(c *fiber.Ctx) error {
		h := adaptor.HTTPHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		}))
		return h(c)
	})

	return adaptor.FiberApp(p.fiberApp)
}

func (p *Plugin) Stop(_ context.Context) error {
	p.log.Info("Closing open resources")

	return nil
}

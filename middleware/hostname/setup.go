package hostname

import (
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/middleware"

	"github.com/mholt/caddy"
)

func init() {
	caddy.RegisterPlugin("hostname", caddy.Plugin{
		ServerType: "dns",
		Action:     setupHostname,
	})
}

func setupHostname(c *caddy.Controller) error {
	c.Next()
	if c.NextArg() {
		return middleware.Error("hostname", c.ArgErr())
	}

	dnsserver.GetConfig(c).AddMiddleware(func(next middleware.Handler) middleware.Handler {
		return Hostname{}
	})

	return nil
}

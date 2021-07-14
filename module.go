package template

import (
	"net/url"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/libdns/powerdns"
)

// Provider lets Caddy read and manipulate DNS records hosted by this DNS provider.
type Provider struct{ *powerdns.Provider }

func init() {
	caddy.RegisterModule(Provider{})
}

// CaddyModule returns the Caddy module information.
func (Provider) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "dns.providers.powerdns",
		New: func() caddy.Module { return &Provider{new(powerdns.Provider)} },
	}
}

// Provision sets up the module. Implements caddy.Provisioner.
func (p *Provider) Provision(ctx caddy.Context) error {
	p.Provider.APIToken = caddy.NewReplacer().ReplaceAll(p.Provider.APIToken, "")
	p.Provider.ServerID = caddy.NewReplacer().ReplaceAll(p.Provider.ServerID, "")
	p.Provider.ServerURL = caddy.NewReplacer().ReplaceAll(p.Provider.ServerURL, "")
	p.Provider.Debug = caddy.NewReplacer().ReplaceAll(p.Provider.Debug, "")
	return nil
}

// UnmarshalCaddyfile sets up the DNS provider from Caddyfile tokens. Syntax:
//
// providername [<api_token>] {
//     api_token <api_token>
// }
//
// **THIS IS JUST AN EXAMPLE AND NEEDS TO BE CUSTOMIZED.**
func (p *Provider) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		if d.NextArg() {
			return d.ArgErr()
		}
		for nesting := d.Nesting(); d.NextBlock(nesting); {
			switch d.Val() {
			case "api_token":
				if p.Provider.APIToken != "" {
					return d.Err("API token already set")
				}
				p.Provider.APIToken = d.Val()
				if d.NextArg() {
					return d.ArgErr()
				}
			case "server_url":
				p.Provider.ServerURL = d.Val()
				if d.NextArg() {
					return d.ArgErr()
				}
			case "server_id":
				p.Provider.ServerID = d.Val()
				if d.NextArg() {
					return d.ArgErr()
				}
			case "debug":
				p.Provider.Debug = d.Val()
				if d.NextArg() {
					return d.ArgErr()
				}
			default:
				return d.Errf("unrecognized subdirective '%s'", d.Val())
			}
		}
	}
	if p.Provider.APIToken == "" {
		return d.Err("missing API token")
	}
	if p.Provider.ServerURL == "" {
		return d.Err("missing server_url")
	}
	if _, err := url.Parse(p.Provider.ServerURL); err != nil {
		return d.Err("invalid server_url")
	}

	return nil
}

// Interface guards
var (
	_ caddyfile.Unmarshaler = (*Provider)(nil)
	_ caddy.Provisioner     = (*Provider)(nil)
)

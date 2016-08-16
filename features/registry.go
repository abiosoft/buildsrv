package features

// PluginType describes a plugin type
type PluginType string

// The different types of plugins.
const (
	DirectivePlugin       PluginType = "directive"
	CaddyfileLoaderPlugin            = "caddyfile_loader"
	ServerPlugin                     = "server"
	DNSProviderPlugin                = "dns_provider"
)

// Plugin represents a Caddy plugin.
type Plugin struct {
	Type        PluginType `json:"type"`
	Name        string     `json:"name"`
	Import      string     `json:"import"`                // i.e. the fully qualified package name
	Description string     `json:"description,omitempty"` // does not end with a period
	DocsURL     string     `json:"docs,omitempty"`        // path-absolute ("/docs/...") used in href attributes
	Default     bool       `json:"default,omitempty"`     // if true, this plugin will be selected by default on the download page
	Required    bool       `json:"required,omitempty"`    // if true, this plugin will always be included in a build
}

// Registry is the list of plugins to show on the download page.
// The order does not matter.
var Registry = Plugins{
	// Server types
	{
		Type:        ServerPlugin,
		Name:        "HTTP",
		Import:      "github.com/mholt/caddy/caddyhttp",
		Description: "HTTP server core; everything most sites need",
		Required:    true,
	},

	// Directives
	{
		Type:        DirectivePlugin,
		Name:        "realip",
		Import:      "github.com/captncraig/caddy-realip",
		Description: "Restore original IP when behind a proxy",
		DocsURL:     "/docs/realip",
	},
	{
		Type:        DirectivePlugin,
		Name:        "git",
		Import:      "github.com/abiosoft/caddy-git",
		Description: "Deploy your site with git push",
		DocsURL:     "/docs/git",
	},
	{
		Type:        DirectivePlugin,
		Name:        "locale",
		Import:      "github.com/simia-tech/caddy-locale",
		Description: "Detect locale of client",
		DocsURL:     "/docs/locale",
	},
	{
		Type:        DirectivePlugin,
		Name:        "minify",
		Import:      "github.com/hacdias/caddy-minify",
		Description: "Minify static assets on-the-fly",
		DocsURL:     "/docs/minify",
	},
	{
		Type:        DirectivePlugin,
		Name:        "ipfilter",
		Import:      "github.com/pyed/ipfilter",
		Description: "Block or allow clients based on IP origin",
		DocsURL:     "/docs/ipfilter",
	},
	{
		Type:        DirectivePlugin,
		Name:        "search",
		Import:      "github.com/pedronasser/caddy-search",
		Description: "Site search engine",
		DocsURL:     "/docs/search",
	},
	// TODO: Waiting for captncraig to update cors to the 0.9 plugin format
	{
		Type:        DirectivePlugin,
		Name:        "cors",
		Import:      "github.com/captncraig/cors/caddy",
		Description: "Easily configure Cross-Origin Resource Sharing",
		DocsURL:     "/docs/cors",
	},
	{
		Type:        DirectivePlugin,
		Name:        "jwt",
		Import:      "github.com/BTBurke/caddy-jwt",
		Description: "Authorization with JSON Web Tokens",
		DocsURL:     "/docs/jwt",
	},
	// TODO. Waiting for pschlump to update jsonp to the 0.9 plugin format
	// {
	// 	Type:        DirectivePlugin,
	// 	Name:        "jsonp",
	// 	Import:      "github.com/pschlump/caddy-jsonp",
	// 	Description: "Wrap JSON responses as JSONP",
	// 	DocsURL:     "/docs/jsonp",
	// },
	{
		Type:        DirectivePlugin,
		Name:        "upload",
		Import:      "blitznote.com/src/caddy.upload",
		Description: "Upload files",
		DocsURL:     "/docs/upload",
	},
	{
		Type:        DirectivePlugin,
		Name:        "filemanager",
		Import:      "github.com/hacdias/caddy-filemanager",
		Description: "Manage files on your server with a GUI",
		DocsURL:     "/docs/filemanager",
	},
	{
		Type:        DirectivePlugin,
		Name:        "hugo",
		Import:      "github.com/hacdias/caddy-hugo",
		Description: "Static site generator with admin interface",
		DocsURL:     "/docs/hugo",
	},
	{
		Type:        DirectivePlugin,
		Name:        "mailout",
		Import:      "github.com/SchumacherFM/mailout",
		Description: "SMTP client with REST API and PGP encryption",
		DocsURL:     "/docs/mailout",
	},
	{
		Type:        DirectivePlugin,
		Name:        "prometheus",
		Import:      "github.com/miekg/caddy-prometheus",
		Description: "Prometheus metrics integration",
		DocsURL:     "/docs/prometheus",
	},
	{
		Type:        DirectivePlugin,
		Name:        "ratelimit",
		Import:      "github.com/xuqingfeng/caddy-rate-limit",
		Description: "Limit rate of requests",
		DocsURL:     "/docs/ratelimit",
	},

	// DNS providers
	{
		Type:   DNSProviderPlugin,
		Name:   "cloudflare",
		Import: "github.com/caddyserver/dnsproviders/cloudflare",
	},
	{
		Type:   DNSProviderPlugin,
		Name:   "digitalocean",
		Import: "github.com/caddyserver/dnsproviders/digitalocean",
	},
	{
		Type:   DNSProviderPlugin,
		Name:   "dnsimple",
		Import: "github.com/caddyserver/dnsproviders/dnsimple",
	},
	{
		Type:   DNSProviderPlugin,
		Name:   "dyn",
		Import: "github.com/caddyserver/dnsproviders/dyn",
	},
	{
		Type:   DNSProviderPlugin,
		Name:   "gandi",
		Import: "github.com/caddyserver/dnsproviders/gandi",
	},
	{
		Type:   DNSProviderPlugin,
		Name:   "googlecloud",
		Import: "github.com/caddyserver/dnsproviders/googlecloud",
	},
	{
		Type:   DNSProviderPlugin,
		Name:   "namecheap",
		Import: "github.com/caddyserver/dnsproviders/namecheap",
	},
	{
		Type:   DNSProviderPlugin,
		Name:   "rfc2136",
		Import: "github.com/caddyserver/dnsproviders/rfc2136",
	},
	{
		Type:   DNSProviderPlugin,
		Name:   "route53",
		Import: "github.com/caddyserver/dnsproviders/route53",
	},
	{
		Type:   DNSProviderPlugin,
		Name:   "vultr",
		Import: "github.com/caddyserver/dnsproviders/vultr",
	},
}

// Plugins is a list of plugins that can determine
// if a name is a member of the set, and it can also
// write out its directives into a list.
type Plugins []Plugin

// Contains determines if name is in the list p.
func (p Plugins) Contains(name string) bool {
	for _, plug := range p {
		if plug.Name == name {
			return true
		}
	}
	return false
}

// String serializes the list of names into a comma-separated string.
func (p Plugins) String() string {
	if len(p) == 0 {
		return ""
	}
	var s string
	for _, plug := range p {
		s += plug.Name + ","
	}
	return s[:len(s)-1] // trim trailing comma
}

// Packages gets the list of packages in p.
func (p Plugins) Packages() []string {
	imports := make([]string, len(p))
	for i, plug := range p {
		imports[i] = plug.Import
	}
	return imports
}

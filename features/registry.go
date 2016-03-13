package features

// Official add-on middleware registry. To register your middleware,
// add it to this list in the proper order.
var Registry = Middlewares{
	// Essential directives that initialize vital configuration settings
	{"root", "", ""},
	{"tls", "", ""},
	{"bind", "", ""},

	// Other directives that don't necessarily create HTTP handlers (services)
	{"startup", "", ""},
	{"shutdown", "", ""},
	{"git", "github.com/abiosoft/caddy-git", "Deploy your site with git push."},

	// Directives that inject handlers (middleware)
	{"log", "", ""},
	{"gzip", "", ""},
	{"errors", "", ""},
	{"ipfilter", "github.com/pyed/ipfilter", "Block or allow clients based on IP origin."},
	{"search", "github.com/pedronasser/caddy-search", "Activates a site search engine"},
	{"header", "", ""},
	{"cors", "github.com/captncraig/cors/caddy", "Enable Cross Origin Resource Sharing"},
	{"rewrite", "", ""},
	{"redir", "", ""},
	{"ext", "", ""},
	{"mime", "", ""},
	{"basicauth", "", ""},
	{"jsonp", "github.com/pschlump/caddy-jsonp", "Wrap regular JSON responses as JSONP"},
	{"internal", "", ""},
	{"proxy", "", ""},
	{"fastcgi", "", ""},
	{"websocket", "", ""},
	{"markdown", "", ""},
	{"templates", "", ""},
	{"browse", "", ""},
	{"hugo", "github.com/hacdias/caddy-hugo", "Powerful and easy static site generator with admin interface."},
	{"mailout", "github.com/SchumacherFM/mailout", "SMTP client with REST API. Send emails with optional PGP encryption"},
	{"prometheus", "github.com/miekg/caddy-prometheus", "Prometheus metrics"},

}

// Middleware is a directive/package pair
type Middleware struct {
	Directive   string `json:"directive"`
	Package     string `json:"package"`
	Description string `json:"description"`
}

// Middlewares is a list of Middleware that can determine
// if a directive is a member of the set, and it can also
// write out its directives into a list.
type Middlewares []Middleware

// Contains determines if directive is in the list m.
func (m Middlewares) Contains(directive string) bool {
	for _, mid := range m {
		if mid.Directive == directive {
			return true
		}
	}
	return false
}

// String serializes the list of directives into a comma-separated string.
func (m Middlewares) String() string {
	if len(m) == 0 {
		return ""
	}
	var s string
	for _, mid := range m {
		s += mid.Directive + ","
	}
	return s[:len(s)-1] // trim trailing comma
}

// Packages gets the list of packages in m.
func (m Middlewares) Packages() []string {
	imports := make([]string, len(m))
	for i, mid := range m {
		imports[i] = mid.Package
	}
	return imports
}

package delegator

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/bobheadxi/res"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"

	// fork of github.com/go-chi/hostrouter with subdomain wildcard support
	"github.com/RTradeLtd/hostrouter"

	"github.com/RTradeLtd/Nexus/config"
	"github.com/RTradeLtd/Nexus/ipfs"
	"github.com/RTradeLtd/Nexus/log"
	"github.com/RTradeLtd/Nexus/network"
	"github.com/RTradeLtd/Nexus/registry"
	"github.com/RTradeLtd/Nexus/temporal"
)

// Engine manages request delegation
type Engine struct {
	l     *zap.SugaredLogger
	reg   *registry.NodeRegistry
	cache *cache

	networks temporal.PrivateNetworks

	timeout   time.Duration
	keyLookup jwt.Keyfunc
	timeFunc  func() time.Time
	version   string
	domain    string
	direct    bool
}

// EngineOpts denotes options for the delegator engine
type EngineOpts struct {
	Version string
	DevMode bool
	Domain  string

	RequestTimeout time.Duration
	JWTKey         []byte
}

// New instantiates a new delegator engine
func New(l *zap.SugaredLogger, opts EngineOpts,
	reg *registry.NodeRegistry, networks temporal.PrivateNetworks) *Engine {

	var timeFunc = time.Now
	if opts.DevMode {
		timeFunc = func() time.Time { return time.Time{} }
	}

	if opts.RequestTimeout == 0 {
		opts.RequestTimeout = 30 * time.Second
	}

	return &Engine{
		l:     l.Named("delegator"),
		reg:   reg,
		cache: newCache(30*time.Minute, 30*time.Minute),

		networks: networks,

		timeout:   opts.RequestTimeout,
		version:   opts.Version,
		keyLookup: func(t *jwt.Token) (interface{}, error) { return opts.JWTKey, nil },
		timeFunc:  timeFunc,
		domain:    opts.Domain,
	}
}

// Run spins up a server that listens for requests and proxies them appropriately
func (e *Engine) Run(ctx context.Context, opts config.Delegator) error {
	var r = chi.NewRouter()

	// mount middleware
	r.Use(
		cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"HEAD", "GET", "POST", "PUT", "PATCH", "DELETE"},
			AllowedHeaders:   []string{"*"},
			AllowCredentials: true,
		}).Handler,
		middleware.RequestID,
		middleware.RealIP,
		log.NewMiddleware(e.l.Named("requests")),
		middleware.Recoverer,
	)

	// register regular HTTP endpoints
	r.HandleFunc("/status", e.Status)
	// set up network endpoints
	if e.domain != "" {
		e.l.Infow("domain configured - registering subdomain routes", "domain", e.domain)
		e.direct = true

		// handle subdomain-based routing
		hr := hostrouter.New()
		hr.Map("*.api."+e.domain, chi.NewRouter().Route("/", func(r chi.Router) {
			r.Use(e.NetworkAndFeatureSubdomainContext)
			r.HandleFunc("/*", e.Redirect)
		}))
		hr.Map("*.gateway."+e.domain, chi.NewRouter().Route("/", func(r chi.Router) {
			r.Use(e.NetworkAndFeatureSubdomainContext)
			r.HandleFunc("/*", e.Redirect)
		}))
		hr.Map("*.swarm."+e.domain, chi.NewRouter().Route("/", func(r chi.Router) {
			r.Use(e.NetworkAndFeatureSubdomainContext)
			r.HandleFunc("/*", e.Redirect)
		}))
		hr.Map("*.status."+e.domain, chi.NewRouter().Route("/", func(r chi.Router) {
			r.Use(e.NetworkAndFeatureSubdomainContext)
			r.HandleFunc("/*", e.NetworkStatus)
		}))
		// mount the host router
		r.Mount("/", hr)
	} else {
		e.l.Infow("no domain configured - subdomain routes not registered")
		e.direct = false

		// use legacy path-based network features
		r.Route(fmt.Sprintf("/network/{%s}", keyNetwork), func(r chi.Router) {
			r.Use(e.NetworkPathContext)
			r.HandleFunc("/status", e.NetworkStatus)
			r.Route(fmt.Sprintf("/{%s}", keyFeature), func(r chi.Router) {
				r.HandleFunc("/*", e.Redirect)
			})
		})
	}

	// set up server
	var srv = &http.Server{
		Handler: r,

		Addr:         opts.Host + ":" + opts.Port,
		WriteTimeout: e.timeout,
		ReadTimeout:  e.timeout,
	}

	// handle shutdown
	go func() {
		for {
			select {
			case <-ctx.Done():
				shutdown, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				defer cancel()
				if err := srv.Shutdown(shutdown); err != nil {
					e.l.Warnw("error encountered during shutdown", "error", err.Error())
				}
			}
		}
	}()

	// go!
	if opts.TLS.CertPath != "" {
		if err := srv.ListenAndServeTLS(opts.TLS.CertPath, opts.TLS.KeyPath); err != nil && err != http.ErrServerClosed {
			e.l.Errorw("error encountered - service stopped", "error", err)
			return err
		}
	} else {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			e.l.Errorw("error encountered - service stopped", "error", err)
			return err
		}
	}

	return nil
}

// NetworkPathContext creates a handler that injects relevant network context into
// all incoming requests through URL parameters
func (e *Engine) NetworkPathContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var id = chi.URLParam(r, string(keyNetwork))
		n, err := e.reg.Get(id)
		if err != nil {
			res.R(w, r, res.ErrNotFound(err.Error()))
			return
		}

		next.ServeHTTP(w, r.WithContext(
			context.WithValue(r.Context(), keyNetwork, &n),
		))
	})
}

// FeaturePathContext creates a handler that injects relevant feature context into
// all incoming requests through URL parameters
func (e *Engine) FeaturePathContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var feature = chi.URLParam(r, string(keyFeature))
		if !validateFeature(feature) {
			res.R(w, r, res.ErrBadRequest(fmt.Sprintf("invalid feature '%s' requested", feature)))
			return
		}
		next.ServeHTTP(w, r.WithContext(
			context.WithValue(r.Context(), keyFeature, feature),
		))
	})
}

// NetworkAndFeatureSubdomainContext creates a handler that injects relevant network context
// into all incoming requests through what is in the subdomain
func (e *Engine) NetworkAndFeatureSubdomainContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			parts   = strings.Split(r.Host, ".")
			id      = parts[0]
			feature = parts[1]
		)

		if !validateFeature(feature) {
			res.R(w, r, res.ErrBadRequest(fmt.Sprintf("invalid feature '%s' requested", feature)))
			return
		}

		n, err := e.reg.Get(id)
		if err != nil {
			res.R(w, r, res.ErrNotFound(err.Error()))
			return
		}

		next.ServeHTTP(w, r.WithContext(
			context.WithValue(
				context.WithValue(r.Context(),
					keyFeature,
					feature),
				keyNetwork,
				&n),
		))
	})
}

// Redirect manages request redirects
func (e *Engine) Redirect(w http.ResponseWriter, r *http.Request) {
	// retrieve network
	n, ok := r.Context().Value(keyNetwork).(*ipfs.NodeInfo)
	if !ok || n == nil {
		res.R(w, r, res.Err(http.StatusText(422), 422))
		return
	}

	// retrieve requested feature
	feature, ok := r.Context().Value(keyFeature).(string)
	if feature == "" {
		res.R(w, r, res.ErrBadRequest("no feature provided"))
		return
	}

	// set target port based on feature
	var port string
	switch feature {
	case "swarm":
		// Swarm access is open to all by default, since it handles authentication
		// on its own
		port = n.Ports.Swarm
	case "api":
		// IPFS network API access requires an authorized user
		user, err := getUserFromJWT(r, e.keyLookup, e.timeFunc)
		if err != nil {
			res.R(w, r, res.ErrUnauthorized(err.Error()))
			return
		}
		entry, err := e.networks.GetNetworkByName(n.NetworkID)
		if err != nil {
			http.Error(w, "failed to find network", http.StatusNotFound)
			return
		}
		var found = false
		for _, authorized := range entry.Users {
			if user == authorized {
				found = true
			}
		}
		if !found {
			res.R(w, r, res.ErrForbidden("user not authorized"))
			return
		}
		// set access rules
		w.Header().Set("Vary", "Origin")
		if entry.APIAllowedOrigin == "" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		} else {
			w.Header().Set("Access-Control-Allow-Origin", entry.APIAllowedOrigin)
		}
		// catch preflights
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		port = n.Ports.API
	case "gateway":
		// Gateway is only open if configured as such
		if entry, err := e.networks.GetNetworkByName(n.NetworkID); err != nil {
			res.R(w, r, res.ErrNotFound("failed to find network"))
			return
		} else if !entry.GatewayPublic {
			res.R(w, r, res.ErrNotFound("failed to find network gateway"))
			return
		}
		port = n.Ports.Gateway
	default:
		res.R(w, r, res.ErrBadRequest(fmt.Sprintf("invalid feature '%s'", feature)))
		return
	}

	// set up target
	var protocol string
	if r.URL.Scheme != "" {
		protocol = r.URL.Scheme + "://"
	} else {
		protocol = "http://"
	}

	var (
		address = fmt.Sprintf("%s:%s", network.Private, port)
		target  = fmt.Sprintf("%s%s%s", protocol, address, r.RequestURI)
	)

	url, err := url.Parse(target)
	if err != nil {
		res.R(w, r, res.ErrInternalServer("could not parse target", err))
		return
	}

	// set up forwarder, retrieving from cache if available, otherwise set up new
	var proxy *httputil.ReverseProxy
	if proxy = e.cache.Get(fmt.Sprintf("%s-%s", n.NetworkID, feature)); proxy == nil {
		proxy = newProxy(feature, url, e.l, e.direct)
		e.cache.Cache(fmt.Sprintf("%s-%s", n.NetworkID, feature), proxy)
	}

	// serve proxy request
	proxy.ServeHTTP(w, r)
}

// Status reports on proxy status
func (e *Engine) Status(w http.ResponseWriter, r *http.Request) {
	res.R(w, r, res.MsgOK("Nexus proxy is online!",
		"version", e.version))
}

// NetworkStatus reports on the status of a network
func (e *Engine) NetworkStatus(w http.ResponseWriter, r *http.Request) {
	n, ok := r.Context().Value(keyNetwork).(*ipfs.NodeInfo)
	if !ok {
		res.R(w, r, res.Err(http.StatusText(422), 422))
		return
	}

	res.R(w, r, res.MsgOK(fmt.Sprintf("found network %s", n.NetworkID),
		"status", "registered"))
}

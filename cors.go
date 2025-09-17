package possum

import (
	"net/http"
	"strconv"
)

type CORSConfig struct {
	AllowOrigin      string   `mapstructure:"allow_origin,omitempty"`
	AllowMethods     string   `mapstructure:"allow_methods,omitempty"`
	AllowHeaders     string   `mapstructure:"allow_headers,omitempty"`
	AllowCredentials bool     `mapstructure:"allow_credentials,omitempty"`
	ExposeHeaders    string   `mapstructure:"expose_headers,omitempty"`
	MaxAge           int      `mapstructure:"max_age,omitempty"`
	ExemptMethods    []string `mapstructure:"exempt_methods,omitempty"`

	cachedMethods map[string]struct{}
}

func (config *CORSConfig) Init() {
	config.cachedMethods = make(map[string]struct{})
	for _, method := range config.ExemptMethods {
		config.cachedMethods[method] = struct{}{}
	}
}

func (config *CORSConfig) SkipMethod(method string) bool {
	_, ok := config.cachedMethods[method]
	return ok
}

var defaultCORSConfig = &CORSConfig{
	AllowOrigin:      "*",
	AllowMethods:     "*",
	AllowHeaders:     "*",
	AllowCredentials: true,
	ExposeHeaders:    "*",
	MaxAge:           0,
	ExemptMethods:    []string{http.MethodOptions},
}

func init() {
	defaultCORSConfig.Init()
}

func SetDefaultCORSConfig(config *CORSConfig) {
	defaultCORSConfig = config
}

// Cors returns a middleware that handles Cross-Origin Resource Sharing (CORS) headers for HTTP requests.
// When used with the Chain function, it should be called as Cors() to return the middleware.
func Cors(config *CORSConfig) HandlerFunc {
	if config == nil {
		config = defaultCORSConfig
	}
	return func(next http.HandlerFunc) http.HandlerFunc {
		return corsHandler(config, next)
	}
}

// CorsHandler is the default CORS middleware that uses the default configuration.
// This can be used directly as CorsHandler(next) or as Cors(nil) when chaining.
func CorsHandler(next http.HandlerFunc) http.HandlerFunc {
	return corsHandler(defaultCORSConfig, next)
}

// corsHandler is the actual CORS middleware implementation.
func corsHandler(config *CORSConfig, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		if config.AllowOrigin != "" {
			origin := r.Header.Get("Origin")
			// When allowing credentials, we must set the specific origin rather than "*" if origin is provided
			if config.AllowCredentials && config.AllowOrigin == "*" && origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			} else {
				w.Header().Set("Access-Control-Allow-Origin", config.AllowOrigin)
			}
			// Add Vary header when we have a specific origin pattern
			if config.AllowOrigin != "*" {
				w.Header().Set("Vary", "Origin")
			}
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "null")
		}

		// Handle credentials
		if config.AllowCredentials {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		// Handle methods
		if config.AllowMethods != "" {
			w.Header().Set("Access-Control-Allow-Methods", config.AllowMethods)
		}

		// Handle headers
		if config.AllowHeaders != "" {
			w.Header().Set("Access-Control-Allow-Headers", config.AllowHeaders)
		}

		// Handle exposed headers
		if config.ExposeHeaders != "" {
			w.Header().Set("Access-Control-Expose-Headers", config.ExposeHeaders)
		}

		// Handle max age
		if config.MaxAge > 0 {
			w.Header().Set("Access-Control-Max-Age", strconv.Itoa(config.MaxAge))
		}

		// Handle WebSocket upgrade requests
		if r.Header.Get("Upgrade") == "websocket" {
			// Ensure WebSocket requests allow necessary WebSocket headers
			w.Header().Set("Access-Control-Allow-Headers",
				config.AllowHeaders+", Sec-WebSocket-Key, Sec-WebSocket-Protocol, Sec-WebSocket-Version")
		}

		// Handle preflight requests
		if config.SkipMethod(r.Method) {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

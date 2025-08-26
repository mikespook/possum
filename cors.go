package possum

import (
	"net/http"
	"strconv"
)

type CORSConfig struct {
	AllowOrigin    string   `mapstructure:"allow_origin,omitempty"`
	AllowedOrigins []string `mapstructure:"allowed_origins,omitempty"`

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

// isOriginAllowed checks if the request origin is allowed based on CORS configuration.
func isOriginAllowed(config *CORSConfig, origin string) bool {
	if origin == "" {
		return false
	}

	// 如果配置了 "*" 且不允许 credentials，直接返回true
	if config.AllowOrigin == "*" && !config.AllowCredentials {
		return true
	}

	// 检查具体的源
	for _, allowed := range config.AllowedOrigins {
		if allowed == origin {
			return true
		}
	}
	return false
}

// Cors is a middleware that handles Cross-Origin Resource Sharing (CORS) headers for HTTP requests.
func Cors(config *CORSConfig, next http.HandlerFunc) http.HandlerFunc {
	if config == nil {
		config = defaultCORSConfig
	}
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// 处理源
		if origin != "" {
			if isOriginAllowed(config, origin) {
				// 当允许credentials时，必须设置具体的源而不是 "*"
				w.Header().Set("Access-Control-Allow-Origin", origin)
			} else {
				// 默认配置下，允许所有源
				w.Header().Set("Access-Control-Allow-Origin", "*")
			}
			w.Header().Set("Vary", "Origin")
		}

		// 处理方法
		if config.AllowMethods != "" {
			w.Header().Set("Access-Control-Allow-Methods", config.AllowMethods)
		}

		// 处理头部
		if config.AllowHeaders != "" {
			w.Header().Set("Access-Control-Allow-Headers", config.AllowHeaders)
		}

		// 处理暴露的头部
		if config.ExposeHeaders != "" {
			w.Header().Set("Access-Control-Expose-Headers", config.ExposeHeaders)
		}

		// 处理缓存时间
		if config.MaxAge > 0 {
			w.Header().Set("Access-Control-Max-Age", strconv.Itoa(config.MaxAge))
		}

		// 处理credentials
		if config.AllowCredentials {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		// 处理WebSocket升级请求
		if r.Header.Get("Upgrade") == "websocket" {
			// 确保WebSocket请求允许的头部包含必要的WebSocket头部
			w.Header().Set("Access-Control-Allow-Headers",
				config.AllowHeaders+", Sec-WebSocket-Key, Sec-WebSocket-Protocol, Sec-WebSocket-Version")
		}

		// 处理预检请求
		if config.SkipMethod(r.Method) {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

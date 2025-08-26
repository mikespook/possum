package possum

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mikespook/possum/config"
)

const (
	// WebSocket configuration constants for timeouts and message sizes.
	writeWait      = 10 * time.Second    // 写超时
	pongWait       = 60 * time.Second    // 等待pong响应超时
	pingPeriod     = (pongWait * 9) / 10 // 发送ping间隔（略小于pongWait）
	maxMessageSize = 512 * 1024          // 最大消息大小（512KB）
)

var (
	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			// TODO: 根据配置检查origin
			return true
		},
	}
)

type WebsocketHandlerFunc func(conn *websocket.Conn, r *http.Request)

// WebSocketUpgrade is a middleware that handles WebSocket connections with CORS support.
func WebSocketUpgrade(corsConfig *CORSConfig, next WebsocketHandlerFunc) http.HandlerFunc {
	if corsConfig == nil {
		corsConfig = defaultCORSConfig
	}
	// 配置WebSocket upgrader
	websocketUpgrader.ReadBufferSize = 1024  // 可以从配置中读取
	websocketUpgrader.WriteBufferSize = 1024 // 可以从配置中读取

	// 根据环境和CORS配置设置CheckOrigin
	if config.IsDev() {
		websocketUpgrader.CheckOrigin = func(r *http.Request) bool {
			return true // 开发环境允许所有源
		}
	} else {
		websocketUpgrader.CheckOrigin = func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			if corsConfig.AllowOrigin == "*" {
				return true // 如果允许所有源，直接返回true
			}
			if len(corsConfig.AllowedOrigins) == 0 {
				return true // 如果没有配置具体的允许源，默认允许
			}
			// 检查origin是否在允许的列表中
			for _, allowed := range corsConfig.AllowedOrigins {
				if allowed == "*" || allowed == origin {
					return true
				}
			}
			return false
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// 升级HTTP连接到WebSocket
		conn, err := websocketUpgrader.Upgrade(w, r, nil)
		if err != nil {
			// WebSocket升级失败时，Upgrade函数已经写入了错误响应，不需要再次写入
			// 避免重复的WriteHeader调用
			return
		}

		// 设置连接参数
		conn.SetReadLimit(maxMessageSize)
		conn.SetReadDeadline(time.Now().Add(pongWait))
		conn.SetPongHandler(func(string) error {
			conn.SetReadDeadline(time.Now().Add(pongWait))
			return nil
		})

		// 启动ping处理器
		go func() {
			ticker := time.NewTicker(pingPeriod)
			defer ticker.Stop()
			for {
				<-ticker.C
				conn.SetWriteDeadline(time.Now().Add(writeWait))
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					return
				}
			}
		}()

		// 确保连接关闭和清理
		defer func() {
			ticker := time.NewTicker(writeWait)
			defer ticker.Stop()
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			conn.Close()
		}()

		// 调用下一个处理器
		next(conn, r)
	}
}

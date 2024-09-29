package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func (m *Middleware) RateLimitMiddleware(rps int, burst int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := fmt.Sprintf("ratelimit:%s", ip)

		allowed, err := m.repos.MiddlewareCache.Allow(key, int64(burst))
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": "Rate limit error"})
			return
		}

		if !allowed {
			c.AbortWithStatusJSON(429, gin.H{"error": "Rate limit exceeded"})
			return
		}

		c.Next()
	}
}

// // RateLimitMiddleware 创建一个 Gin 中间件来实现 IP 基础的速率限制
// // rps: 每秒允许的请求数
// // burst: 允许的突发请求数
// func (m *Middleware) RateLimitMiddleware(rps int, burst int) gin.HandlerFunc {
// 	// client 结构体用于存储每个 IP 的限流器和最后访问时间
// 	type client struct {
// 		limiter  *rate.Limiter
// 		lastSeen time.Time
// 	}

// 	var (
// 		// 使用 map 存储每个 IP 的 client 信息
// 		clients = make(map[string]*client)
// 		// 互斥锁用于保护 clients map 的并发访问
// 		mu sync.Mutex
// 	)

// 	// 启动一个后台 goroutine 来清理长时间未活动的 IP
// 	go func() {
// 		for {
// 			time.Sleep(time.Minute)
// 			mu.Lock()
// 			for ip, client := range clients {
// 				// 如果 IP 超过 3 分钟没有活动，则从 map 中删除
// 				if time.Since(client.lastSeen) > 3*time.Minute {
// 					delete(clients, ip)
// 				}
// 			}
// 			mu.Unlock()
// 		}
// 	}()

// 	// 返回实际的中间件处理函数
// 	return func(c *gin.Context) {
// 		// 获取客户端 IP
// 		ip := c.ClientIP()

// 		mu.Lock()
// 		// 如果该 IP 不存在，则创建一个新的限流器
// 		if _, exists := clients[ip]; !exists {
// 			clients[ip] = &client{limiter: rate.NewLimiter(rate.Limit(rps), burst)}
// 		}
// 		// 更新最后访问时间
// 		clients[ip].lastSeen = time.Now()

// 		// 检查是否允许本次请求
// 		if !clients[ip].limiter.Allow() {
// 			mu.Unlock()
// 			// 如果超出速率限制，返回 429 状态码
// 			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
// 			return
// 		}
// 		mu.Unlock()

// 		// 继续处理请求
// 		c.Next()
// 	}
// }

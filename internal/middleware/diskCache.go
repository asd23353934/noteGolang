package middleware

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	cacheDir = "./cache"
)

type CacheConfig struct {
	Duration time.Duration
}

func (m *Middleware) DiskCache(config CacheConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := getCacheKey(c.Request)
		cacheFile := filepath.Join(cacheDir, key)

		// 檢查快取是否存在且未過期
		if cacheExists(cacheFile, config.Duration) {
			serveJSONCache(c, cacheFile)
			return
		}

		// 創建一個 writer 來捕獲響應
		writer := &jsonResponseWriter{ResponseWriter: c.Writer}
		c.Writer = writer

		// 處理請求
		c.Next()

		// 如果響應成功或是 304 Not Modified，則快取
		if c.Writer.Status() == http.StatusOK || c.Writer.Status() == http.StatusNotModified {
			saveJSONCache(cacheFile, writer.body)
		}
	}
}

func getCacheKey(r *http.Request) string {
	hash := md5.Sum([]byte(r.URL.RequestURI()))
	return hex.EncodeToString(hash[:])
}

func cacheExists(filename string, duration time.Duration) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return time.Since(info.ModTime()) < duration
}

func serveJSONCache(c *gin.Context, filename string) {
	info, err := os.Stat(filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read cache"})
		return
	}

	// 計算 ETag
	etag := fmt.Sprintf(`"%x"`, md5.Sum([]byte(fmt.Sprintf("%v", info.ModTime().Unix()))))

	// 處理條件請求
	if match := c.GetHeader("If-None-Match"); match != "" {
		if match == etag {
			c.Status(http.StatusNotModified)
			return
		}
	}

	if modifiedSince := c.GetHeader("If-Modified-Since"); modifiedSince != "" {
		if t, err := time.Parse(http.TimeFormat, modifiedSince); err == nil {
			if info.ModTime().Unix() <= t.Unix() {
				c.Status(http.StatusNotModified)
				return
			}
		}
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read cache"})
		return
	}

	// 設置快取控制頭
	c.Header("Cache-Control", "public, max-age=300") // 可以根據需求調整 max-age
	c.Header("ETag", etag)
	c.Header("Last-Modified", info.ModTime().UTC().Format(http.TimeFormat))
	c.Header("X-Cache", "HIT")
	c.Header("Content-Type", "application/json")

	// 設置 Age 頭
	age := time.Since(info.ModTime()).Seconds()
	c.Header("Age", fmt.Sprintf("%.0f", age))

	c.Writer.Write(data)
	c.Abort()
}

func saveJSONCache(filename string, content []byte) {
	os.MkdirAll(filepath.Dir(filename), 0755)
	os.WriteFile(filename, content, 0644)
}

type jsonResponseWriter struct {
	gin.ResponseWriter
	body []byte
}

func (w *jsonResponseWriter) Write(b []byte) (int, error) {
	w.body = append(w.body, b...)
	return w.ResponseWriter.Write(b)
}

func (w *jsonResponseWriter) WriteJSON(data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.body = jsonData
	_, err = w.ResponseWriter.Write(jsonData)
	return err
}

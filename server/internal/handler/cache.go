package handler

import (
	"encoding/json"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/klkl2164/urlvet/internal/service/cache"
)

// cacheEntry is the per-key payload returned by ListCacheHandler.
type cacheEntry struct {
	Key        string          `json:"key"`
	Prefix     string          `json:"prefix"`
	TTLSeconds int64           `json:"ttl_seconds"` // -1 = no expiry, -2 = key gone
	Value      json.RawMessage `json:"value" swaggertype:"object"`
}

// ListCacheHandler godoc
// @Summary      List all cache keys and values
// @Tags         admin
// @Produce      json
// @Success      200  {object}  map[string]any
// @Router       /admin/cache [get]
func ListCacheHandler(c *gin.Context) {
	cacheInstance, err := cache.New()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法连接到缓存服务"})
		return
	}
	defer cacheInstance.Close()

	ctx := c.Request.Context()

	keys, err := cacheInstance.Scan(ctx, "*")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法列出缓存键"})
		return
	}
	sort.Strings(keys)

	entries := make([]cacheEntry, 0, len(keys))
	for _, key := range keys {
		raw, err := cacheInstance.Get(ctx, key)
		if err != nil || raw == "" {
			continue
		}

		ttl, _ := cacheInstance.TTL(ctx, key)
		ttlSec := int64(ttl.Seconds())

		// Try to present value as formatted JSON; fall back to a plain string.
		var valueRaw json.RawMessage
		if json.Valid([]byte(raw)) {
			valueRaw = json.RawMessage(raw)
		} else {
			quoted, _ := json.Marshal(raw)
			valueRaw = json.RawMessage(quoted)
		}

		prefix := key
		for i, ch := range key {
			if ch == ':' {
				prefix = key[:i]
				break
			}
		}

		entries = append(entries, cacheEntry{
			Key:        key,
			Prefix:     prefix,
			TTLSeconds: ttlSec,
			Value:      valueRaw,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"total": len(entries),
		"keys":  entries,
	})
}

// DeleteCacheKeyHandler godoc
// @Summary      Delete a single cache key
// @Tags         admin
// @Param        key  path  string  true  "Cache key"
// @Success      200  {object}  map[string]any
// @Router       /admin/cache/{key} [delete]
func DeleteCacheKeyHandler(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少缓存键"})
		return
	}
	// Gin includes the leading slash when using wildcard params; strip it.
	if len(key) > 0 && key[0] == '/' {
		key = key[1:]
	}

	cacheInstance, err := cache.New()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法连接到缓存服务"})
		return
	}
	defer cacheInstance.Close()

	if err := cacheInstance.Delete(c.Request.Context(), key); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法删除缓存键"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted", "key": key})
}

// FlushCacheHandler godoc
// @Summary      Flush all cache keys
// @Tags         admin
// @Success      200  {object}  map[string]any
// @Router       /admin/cache [delete]
func FlushCacheHandler(c *gin.Context) {
	cacheInstance, err := cache.New()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "ERROR",
			"error":  "无法连接到缓存服务",
		})
		return
	}
	defer cacheInstance.Close()

	if err := cacheInstance.FlushAll(c.Request.Context()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "ERROR",
			"error":  "无法清空缓存",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "SUCCESS",
		"message": "缓存已清空",
	})
}

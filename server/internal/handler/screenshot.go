package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/klkl2164/urlvet/internal/logger"
	"github.com/klkl2164/urlvet/internal/service/checks"
	"github.com/klkl2164/urlvet/internal/service/screenshot"
)

func ScreenshotHandler(c *gin.Context) {
	rawURL := c.Query("url")
	if rawURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少链接参数"})
		return
	}

	// Basic URL validation (additional validation happens in TakeScreenshot)
	_, isValid, err := checks.IsValidURL(rawURL)
	if err != nil || !isValid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "链接格式不正确"})
		return
	}

	imageBytes, err := screenshot.TakeScreenshot(rawURL)
	if err != nil {
		// The error can name the target URL, so only its occurrence is logged.
		logger.Error("screenshot failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "截图失败"})
		return
	}

	// Return the image as binary with proper content-type
	c.Data(http.StatusOK, "image/png", imageBytes)
}

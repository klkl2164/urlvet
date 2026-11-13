package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/klkl2164/urlvet/internal/service/checks"
)

// CheckHSTSHandler checks whether the URL's host enforces HSTS.
//
//	@Summary		HSTS check
//	@Tags			Security
//	@Produce		json
//	@Param			url	query		string	true	"URL to check"
//	@Success		200	{object}	map[string]bool
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/hsts [get]
func CheckHSTSHandler(c *gin.Context) {
	rawURL := c.Query("url")
	if rawURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少链接参数"})
		return
	}

	_, isValid, err := checks.IsValidURL(rawURL)
	if err != nil || !isValid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "链接格式不正确"})
		return
	}

	supportsHSTS, err := checks.SupportsHSTS(rawURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "检查 HSTS 失败（" + err.Error() + "）"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"supports_hsts": supportsHSTS,
	})
}

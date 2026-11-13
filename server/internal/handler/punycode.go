package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/klkl2164/urlvet/internal/service/checks"
)

// CheckPunycodeHandler detects punycode (IDN) characters in the URL.
//
//	@Summary		Punycode detection
//	@Tags			URL
//	@Produce		json
//	@Param			url	query		string	true	"URL to check"
//	@Success		200	{object}	map[string]bool
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/punycode [get]
func CheckPunycodeHandler(c *gin.Context) {
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

	hasPunycode, err := checks.ContainsPunycode(rawURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "检查 Punycode 失败（" + err.Error() + "）"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"contains_punycode": hasPunycode,
	})
}

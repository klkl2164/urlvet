package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/klkl2164/urlvet/internal/service/checks"
	"github.com/klkl2164/urlvet/internal/service/domaininfo"
)

// DomainInfoHandler returns WHOIS / RDAP registration data for the URL's domain.
//
//	@Summary		Domain WHOIS / RDAP info
//	@Tags			DNS
//	@Produce		json
//	@Param			url	query		string	true	"URL whose domain to look up"
//	@Success		200	{object}	map[string]any
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/domain-info [get]
func DomainInfoHandler(c *gin.Context) {
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

	domain, err := checks.GetDomain(rawURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无法从链接中提取域名"})
		return
	}

	domainInfo, err := domaininfo.Lookup(domain)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "域名信息查询失败",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"domain":     domain,
		"domainInfo": domainInfo,
	})
}

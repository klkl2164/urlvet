// This handler is for development/testing purposes only.
// This code is intended for local experiments and feature prototyping.
// Do NOT commit permanent changes or production features to this handler.

package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/klkl2164/urlvet/internal/service/checks"
)

func TestApiHandler(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少链接参数"})
		return
	}

	_, isValid, err := checks.IsValidURL(url)
	if err != nil || !isValid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "链接格式不正确"})
		return
	}

	domain, err := checks.GetDomain(url)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无法从链接中提取域名"})
		return
	}

	////////////////////////////////////////////////////////////////////
	// THIS IS A TEMPORARY CODE SECTION
	// MAKE CHANGES ONLY INSIDE THIS BOX
	// SAFE TO DELETE LATER IF PUSHED BY MISTAKE

	c.JSON(http.StatusOK, gin.H{
		"domain": domain,
	})

	////////////////////////////////////////////////////////////////////

}

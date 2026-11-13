package handler

import (
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/klkl2164/urlvet/internal/admintoken"
	"golang.org/x/crypto/argon2"
)

type adminLoginRequest struct {
	Password string `json:"password" binding:"required"`
}

// AdminLoginHandler validates the admin password against the stored Argon2id hash
// and returns a signed session token on success.
func AdminLoginHandler(c *gin.Context) {
	hash := os.Getenv("ADMIN_PASSWORD_HASH")
	if hash == "" {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "管理后台没有启用"})
		return
	}

	var req adminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请输入密码"})
		return
	}

	if err := verifyArgon2id(req.Password, hash); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "密码不正确"})
		return
	}

	token, err := admintoken.Issue(24 * time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法生成登录令牌"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// verifyArgon2id checks a password against a PHC-format Argon2id hash:
// $argon2id$v=19$m=65536,t=3,p=4$<base64_salt>$<base64_hash>
func verifyArgon2id(password, encodedHash string) error {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		return errors.New("invalid hash format")
	}

	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil || version != argon2.Version {
		return errors.New("unsupported argon2 version")
	}

	var m, t, p int
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &m, &t, &p); err != nil {
		return errors.New("invalid hash parameters")
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return errors.New("invalid salt encoding")
	}

	storedKey, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return errors.New("invalid hash encoding")
	}

	computedKey := argon2.IDKey([]byte(password), salt, uint32(t), uint32(m), uint8(p), uint32(len(storedKey)))
	if subtle.ConstantTimeCompare(computedKey, storedKey) != 1 {
		return errors.New("password mismatch")
	}
	return nil
}

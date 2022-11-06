package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/IOPaper/Paper/utils/result"
	"github.com/gin-gonic/gin"
)

func AdminAuth(secret string) gin.HandlerFunc {
	var hashSecret string
	{
		s := sha256.New()
		s.Write([]byte(secret))
		hashSecret = hex.EncodeToString(s.Sum(nil))
	}
	return func(c *gin.Context) {
		if c.GetHeader("SECRET-v1") != hashSecret {
			result.New[any]().SetStatusCode(401).SetMessage("unauthorized").Err(c.Writer)
			c.Abort()
			return
		}
		c.Next()
	}
}

//go:build module_auth

package module

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"github.com/IOPaper/Paper/global"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Simple module demo, authentication

type WebAuth struct {
	Token string
}

func NewWebAuth() Call {
	return &WebAuth{Token: global.Config.Engine.AuthToken}
}

func (a *WebAuth) verify(cc, cr string) error {
	c, err := base64.URLEncoding.DecodeString(cr)
	if err != nil {
		return errors.New("invalid challenge")
	}
	h := hmac.New(sha256.New, []byte(a.Token))
	h.Write([]byte(cc))
	result := h.Sum(nil)
	if !hmac.Equal(result, c) {
		return errors.New("invalid challenge")
	}
	return nil
}

func (a *WebAuth) Action() gin.HandlerFunc {
	return func(c *gin.Context) {
		cc := c.GetHeader("challenge-content")
		cr := c.GetHeader("challenge-result")
		if cc == "" || cr == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if err := a.verify(cc, cr); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		c.Next()
	}
}

func init() {
	registerCall(NewWebAuth())
}

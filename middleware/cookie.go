package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/eirka/eirka-libs/config"
	e "github.com/eirka/eirka-libs/errors"
)

func GetAntiSpamCookie() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Test for cookie from Prim
		cookie, err := c.Request.Cookie(config.Settings.Antispam.CookieName)
		if err == http.ErrNoCookie {
			c.JSON(http.StatusBadRequest, gin.H{"error_message": e.ErrNoCookie.Error()})
			c.Error(e.ErrNoCookie)
			c.Abort()
			return
		}

		// See if the cookie is the right value
		if cookie.Value != config.Settings.Antispam.CookieValue {
			c.JSON(http.StatusBadRequest, gin.H{"error_message": e.ErrInvalidCookie.Error()})
			c.Error(e.ErrInvalidCookie)
			c.Abort()
			return
		}

		c.Next()

	}
}

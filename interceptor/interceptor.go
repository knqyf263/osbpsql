package interceptor

import (
	"github.com/knqyf263/osbpsql/config"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// BasicAuth authenticates Basic authentication
func BasicAuth() echo.MiddlewareFunc {
	return middleware.BasicAuth(func(username string, password string, context echo.Context) (bool, error) {
		basicUsername := config.Config.BasicAuthUsername
		basicPassword := config.Config.BasicAuthPassword
		if username == basicUsername && password == basicPassword {
			return true, nil
		}
		return false, nil
	})
}

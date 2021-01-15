package logger
import "go.uber.org/zap"
import echo "github.com/labstack/echo/v4"
import "net/http"
//import "github.com/labstack/echo/v4/middleware"

const loggerKey = "logger"

func MiddlewareLogger(logger *zap.Logger) echo.MiddlewareFunc{
	
	
	return func(next echo.HandlerFunc) echo.HandlerFunc{
		return func(c echo.Context) error{
			
			id := c.Request().Header.Get("X-Request-ID")
			if id == ""{
				return echo.NewHTTPError(http.StatusUnauthorized, "Please provide X-Request-ID header")
			}

			l:= logger.With(zap.String("x-request-id", id))
			c.Set(loggerKey,l)			
			
			//return nil
			err := next(c)

			return err
		}
	}
}

func Extract(c echo.Context) *zap.Logger{
	l,ok := c.Get(loggerKey).(*zap.Logger)
	if ok {
		return l
	}

	return zap.NewExample()
}
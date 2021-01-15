package auth

import (
	"encoding/json"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)
import echo "github.com/labstack/echo/v4"
import "github.com/pallat/todos/logger"
import "errors"
import "strings"

func TokenHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	t, err := GenerateToken()
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"token": t,
	})
}

var Token = GenerateToken

// token expire 15 minutes
func GenerateToken() (string, error) {
	mySigningKey := []byte("drowssap")

	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Minute*15).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(mySigningKey)
}

func ValidateToken(tokenString string,c echo.Context) (jwt.Claims,error){
	logger :=  logger.Extract(c)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("drowssap"), nil
	})

	//logger.Info(" ValidateToken after parse ")
	

	if err == nil && token.Valid {
		logger.Info("token valid")
		return token.Claims,nil
		//logger.Info("You look nice today")
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			logger.Info("invalid token(1)")
			return nil,errors.New("invalid token(1)")
			//logger.Info("That's not even a token")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			logger.Info("token expire")
			return nil,errors.New("token expire")
			//logger.Info("Timing is everything")
		} else {
			logger.Info("invalid token(2)")
			return nil, errors.New("invalid token(2)")
			//logger.Info("Couldn't handle this token:", err)
		}
	} else {
		logger.Info("invalid token(3)")
		return nil, errors.New("invalid token(3)")
		//logger.Info("Couldn't handle this token:", err)
	}
}

func MiddlewareJwtAuthen() echo.MiddlewareFunc{
	
	
	return func(next echo.HandlerFunc) echo.HandlerFunc{
		return func(c echo.Context) error{
			
			authorization := c.Request().Header.Get("Authorization")
			
			if authorization == ""{
				// c.Response().Header().Set(echo.HeaderServer, "Echo/3.0")
				// return error.New("")
				return echo.NewHTTPError(http.StatusUnauthorized, "Please provide valid credentials")
				//c.String(http.StatusUnauthorized)
			}

			logger :=  logger.Extract(c)


			token := strings.ReplaceAll(authorization, "Bearer ", "")

			logger.Info("token : "+ token);

			claims,errorValidate := ValidateToken(token,c)

			_ = claims
			if errorValidate != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Please provide valid credentials")
			}
			
			//return nil
			
			err := next(c)

			return err
		}
	}
}
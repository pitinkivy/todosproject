package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"strings"

	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/spf13/viper"
	"github.com/joho/godotenv"
	"github.com/pallat/todos/captcha"
	"github.com/pallat/todos/auth"
	"github.com/pallat/todos/todos"
	"github.com/pallat/todos/logger"
	"go.uber.org/zap"
)

import (
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
  )

func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")
  
	if err != nil {
	  log.Fatalf("Error loading .env file")
	}
  
	return os.Getenv(key)
}

func main() {
	//goDotEnvVariable()

	


	viper.SetDefault("app.addr", "0.0.0.0:8888")
	
	viper.AddConfigPath(".")
	//viper.SetConfigFile(".env")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".","_"))
	
	
	if err:=viper.ReadInConfig();err != nil{
		log.Println("Warning:",err);
	}
	
	l, _ := zap.NewProduction()
	defer l.Sync() // flushes buffer, if any


	// github.com/denisenkom/go-mssqldb
	dsn := viper.GetString("app.db_connection")
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if(err != nil){
		log.Fatal(err);
		/*
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"err":err.Error(),
		})
		*/
	}

	router := echo.New()
	router.Use(middleware.Logger())
	router.Use(middleware.Recover())
	router.Use(logger.MiddlewareLogger(l))
	//router.Use(cors.Default())

	//#region public path
	router.GET("/ping", func(c echo.Context) error {
		return c.String(200, "pong")
	})

	// for get captcha
	router.GET("/captcha", captchaHandler)
	// for exchange captcha answer for jwt token
	router.POST("/exchange", exchangeHandler)
	//#endregion

	//#region secure path using Bearer
	router.POST("/todos", todos.NewTodoHandler(db))
	router.GET("/todos", todos.GetTodoHandler(db))
	router.PUT("/todos/:id", todos.PutUpdateTodoHandler(db))
	router.DELETE("/todos/:id", todos.DeleteTodoHandler(db))
	
	g := router.Group("/todos")
	//_ = g
	g.Use(auth.MiddlewareJwtAuthen())
	//#endregion


	srv := &http.Server{
		Addr:         viper.GetString("app.addr"),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("serve on: %s\n",viper.GetString("app.addr"))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)


	// (syscall.SIGINT)ctrl + c or (syscall.SIGTERM)kill
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// wait 5 second
	/*
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	*/
	
	// wait until job done
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}


func captchaHandler (c echo.Context) error{
	key, captcha := captcha.KeyQuestion();
	return c.JSON(http.StatusOK, map[string]string{
		"key":key,
		"captcha":captcha,
	})
}

func exchangeHandler(c echo.Context) error{
	type Answer struct {
		Key string `json:"key"`
		Answer int `json:"answer"`
	}

	var ans Answer
	if err := c.Bind(&ans); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":err.Error(),
		})
	}

	if !captcha.Answer(ans.Key, ans.Answer){
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error":"wrong answer",
		})
	}

	t, err :=  auth.Token()
	if err != nil{
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token":t,
	})


	/*
	key, captcha := captcha.KeyQuestion();
	return c.JSON(http.StatusOK, map[string]string{
		"key":key,
		"captcha":captcha,
	})*/
}

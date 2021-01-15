package todos

import "fmt"
import echo "github.com/labstack/echo/v4"
import "net/http"
import "log"
import "strconv"
import goerror "errors"
//import "go.uber.org/zap"
import "github.com/pallat/todos/logger"
import "github.com/pkg/errors"
import "go.uber.org/zap"
import (
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
  )

  type Todo struct{
	Task string `json:"task"`
}

// POST /todos for create todo record
func NewTodoHandler(db *gorm.DB) echo.HandlerFunc{
	//return newTodoHandlerLocal;
	return func(c echo.Context) error {
		var todo Todo

		var funcAlias = fmt.Sprintf("POST /todos")

		// assert type
		logger :=  logger.Extract(c)//c.Get("logger").(*zap.Logger)
		logger.Info(funcAlias + " new task todo........")
		if err := c.Bind(&todo); err != nil{
			
			return c.JSON(http.StatusBadRequest , map[string]string{
				"error": errors.Wrap(err,"new task").Error(),
			})
		}
	
		if todo.Task == ""{
			logger.Info(funcAlias+" Task is required ")
			return c.JSON(http.StatusBadRequest , map[string]string{
				"error": "Task is required",
			})
		}
		
	
		if err := db.Create(&Task{
			Task : todo.Task,
		}).Error; err != nil{
			return c.JSON(http.StatusBadRequest, map[string]string{
				"err":errors.Wrap(err,"create task").Error(),
			})
		}
	
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success":true,
		})
		
	}
}

// GET /todos for get all todo records
func GetTodoHandler(db *gorm.DB) echo.HandlerFunc{
	//return newTodoHandlerLocal;
	return func(c echo.Context) error {
		//var todo Todo

		// assert type
		logger :=  logger.Extract(c)//c.Get("logger").(*zap.Logger)
		logger.Info("GET /todos get all todos........")
		
		var tasks = make([]Task,0);
		result := db.Find(&tasks)
		_ = result
	
		return c.JSON(http.StatusOK, tasks)
	}
}

// PUT /todos/:id for update todo record by id
func PutUpdateTodoHandler(db *gorm.DB) echo.HandlerFunc{
	//return newTodoHandlerLocal;
	return func(c echo.Context) error {
		var todo Todo

		logger :=  logger.Extract(c)//c.Get("logger").(*zap.Logger)

		id := c.Param("id")
		var funcAlias = fmt.Sprintf("PUT /todos/%s",id)

		logger.Info(funcAlias+" update id",zap.String("id",id))

		
		if id == ""{
			logger.Info(funcAlias+" id is empty ")
			return c.JSON(http.StatusBadRequest , map[string]string{
				"error":"id is empty" ,
			})
		}

		idInt,err := strconv.Atoi(id)
		if err != nil{
			logger.Info(funcAlias+" id invalid format ",zap.String("id",id))
			return c.JSON(http.StatusBadRequest , map[string]string{
				"error": errors.Wrap(err," id invalid format").Error(),
			})
		}

		// assert type
		
		if err := c.Bind(&todo); err != nil{
			logger.Info(funcAlias+" invalid body ")
			return c.JSON(http.StatusBadRequest , map[string]string{
				"error": errors.Wrap(err,"invalid body").Error(),
			})
		}

		if todo.Task == ""{
			logger.Info(funcAlias+" Task is required ")
			return c.JSON(http.StatusBadRequest , map[string]string{
				"error": "Task is required",
			})
		}
	
		var task Task;
		var result = db.First(&task, idInt)

		if goerror.Is(result.Error, gorm.ErrRecordNotFound){
			logger.Info(funcAlias+" not found record with id ",zap.String("id",id))
			return c.JSON(http.StatusBadRequest , map[string]string{
				"error": "not found record with id "+id,
			})
		}

		task.Task = todo.Task
		result = db.Save(&task)
		//result.RowsAffected // returns updated records count
		//result.Error        // returns updating error
		if(result.RowsAffected == 0){
			logger.Info(funcAlias+" can't update record with id ",zap.String("id",id))
			return c.JSON(http.StatusInternalServerError , map[string]string{
				"error": " can't update record with id "+id,
			})
		}

		return c.JSON(http.StatusOK, task)
	}
}

// DELETE /todos/:id for update todo record by id
func DeleteTodoHandler(db *gorm.DB) echo.HandlerFunc{
	//return newTodoHandlerLocal;
	return func(c echo.Context) error {
		var todo Todo

		logger :=  logger.Extract(c)//c.Get("logger").(*zap.Logger)

		id := c.Param("id")
		var funcAlias = fmt.Sprintf("DELETE /todos/%s",id)

		logger.Info(funcAlias+" delete id",zap.String("id",id))

		
		if id == ""{
			logger.Info(funcAlias+" id is empty ")
			return c.JSON(http.StatusBadRequest , map[string]string{
				"error":"id is empty" ,
			})
		}

		idInt,err := strconv.Atoi(id)
		if err != nil{
			logger.Info(funcAlias+" id invalid format ",zap.String("id",id))
			return c.JSON(http.StatusBadRequest , map[string]string{
				"error": errors.Wrap(err," id invalid format").Error(),
			})
		}

		
		var task Task;
		
		var result = db.Delete(&Task{},idInt)

		
		if goerror.Is(result.Error, gorm.ErrRecordNotFound){
			logger.Info(funcAlias+" not found record with id ",zap.String("id",id))
			return c.JSON(http.StatusBadRequest , map[string]string{
				"error": "not found record with id "+id,
			})
		}

		task.Task = todo.Task
		//result = db.Save(&task)
		//result.RowsAffected // returns updated records count
		//result.Error        // returns updating error
		if(result.RowsAffected == 0){
			logger.Info(funcAlias+" can't delete record with id ",zap.String("id",id))
			return c.JSON(http.StatusInternalServerError , map[string]string{
				"error": " can't delete record with id "+id,
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"success":true,
		})
	}
}

func _newTodoHandlerLocal(c echo.Context) error {
	var todo struct{
		Task string `json:"task"`
	}

	if err := c.Bind(&todo); err != nil{
		return c.JSON(http.StatusBadRequest , map[string]string{
			"error":err.Error(),
		})
	}

	dsn := "sqlserver://test:555@localhost:1434?database=go_training"
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if(err != nil){
		log.Fatal(err);
		/*
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"err":err.Error(),
		})
		*/
	}

	
	
	//db.AutoMigrate(Task{})

	//var todoTable Task
	//user := User{Name: "Jinzhu", Age: 18, Birthday: time.Now()}

	/*result := db.Create(&todoTable) // pass pointer of data to Create

	if result.Error != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"err":result.Error.Error(),
		})
	}*/

	if err := db.Create(&Task{
		Task : todo.Task,
	}).Error; err != nil{
		return c.JSON(http.StatusBadRequest, map[string]string{
			"err":err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{})
}

type Task struct{
	gorm.Model
	Task string
	Processed bool
}

func (Task) TableName() string{
	return "todos"
}
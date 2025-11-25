package handlers

import (
	"net/http"

	"github.com/ValianceTekProject/AreaBack/db"
	"github.com/ValianceTekProject/AreaBack/initializers"
	"github.com/gin-gonic/gin"
)

type todo struct {
	Item      string `json:"Item"`
	Completed bool   `json:"completed"`
}

func PostTodos(context *gin.Context) {
	var newTodo todo
	if err := context.BindJSON(&newTodo); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	_, err := initializers.DB.Todo.CreateOne(
		db.Todo.Item.Set(newTodo.Item),
		db.Todo.Completed.Set(newTodo.Completed),
	).Exec(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Database insert error"})
		return
	}

	context.Status(http.StatusCreated)
}

func DelTodos(context *gin.Context) {
	var newTodo todo
	if err := context.BindJSON(&newTodo); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	_, err := initializers.DB.Todo.FindUnique(
		db.Todo.Item.Equals(newTodo.Item),
	).Delete().Exec(context)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Database remove error"})
		return
	}

	context.Status(http.StatusOK)
}

func GetTodos(context *gin.Context) {
	todos, err := initializers.DB.Todo.FindMany().Exec(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	context.JSON(http.StatusOK, todos)
}

func ModifyStatus(context *gin.Context) {
	var newTodo todo
	if err := context.BindJSON(&newTodo); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	_, err := initializers.DB.Todo.FindUnique(
		db.Todo.Item.Equals(newTodo.Item),
	).Update(db.Todo.Completed.Set(newTodo.Completed)).Exec(context)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update in database"})
		return
	}
	context.Status(http.StatusOK)
}

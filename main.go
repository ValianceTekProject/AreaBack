package main

import (
	"fmt"

	"net/http"

	"github.com/ValianceTekProject/AreaBack/initializers"
	"github.com/gin-gonic/gin"
)

type todo struct {
	ID			string	`json:"id"`
	Item		string	`json:"Item"`
	Completed	bool	`json:"completed"`
}

var todos = []todo{
	{ID: "1", Item: "clean Room", Completed: false},
	{ID: "2", Item: "Read Book", Completed: false},
	{ID: "3", Item: "Record Video", Completed: false},
}

func getTodos(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, todos)
}

func main() {
    fmt.Println("Hello")
	initializers.ConnectDB();
	router := gin.Default()
	router.GET("/todos", getTodos)
	router.Run("0.0.0.0:8080")
}


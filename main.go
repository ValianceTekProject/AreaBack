package main

import (
	"github.com/ValianceTekProject/AreaBack/handlers"
	"github.com/ValianceTekProject/AreaBack/initializers"
	"github.com/gin-gonic/gin"
)

func main() {
	initializers.ConnectDB()
	router := gin.Default()
	router.POST("/todos/add", handlers.PostTodos)
	router.DELETE("/todos/del", handlers.DelTodos)
	router.Run("0.0.0.0:8080")
}

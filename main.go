package main

import (
	"github.com/ValianceTekProject/AreaBack/initializers"
	"github.com/ValianceTekProject/AreaBack/router"
)

func main() {
	initializers.ConnectDB()
	router.SetupRouting()
}

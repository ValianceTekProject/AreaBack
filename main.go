package main

import (
	"github.com/ValianceTekProject/AreaBack/initializers"
	"github.com/ValianceTekProject/AreaBack/router"
	"github.com/ValianceTekProject/AreaBack/routine"
)

func main() {
	initializers.ConnectDB()
	routine.LaunchRoutines()
	router.SetupRouting()
}

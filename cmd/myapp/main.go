package main

import (
	_ "awesomeProject/internal/api"
	"awesomeProject/internal/pkg/app"
	"log"
)

func main() {
	log.Println("Application start up!")
	a := app.New()
	log.Println("Application created")
	a.StartServer()
	log.Println("Application terminated!")
}

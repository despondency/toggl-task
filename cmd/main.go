package main

import (
	"github.com/despondency/toggl-task/internal/application"
	"log"
)

func main() {
	app := application.NewApplication(8080)
	log.Fatal(app.StartServer())
}

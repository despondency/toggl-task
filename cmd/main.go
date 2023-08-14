package main

import (
	"github.com/despondency/toggl-task/internal/application"
	"log"
)

func main() {
	log.Fatal(application.StartServer(8080))
}

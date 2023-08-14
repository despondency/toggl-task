package application

import (
	"context"
	"fmt"
	v1 "github.com/despondency/toggl-task/internal/handler/v1"
	"github.com/despondency/toggl-task/internal/persister"
	"github.com/despondency/toggl-task/internal/scanner"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type Application struct {
	port int
}

func StartServer(port int) error {
	app := fiber.New(
		fiber.Config{BodyLimit: 4 * 1024 * 1024}, //
	)

	v1Handlers := app.Group("/v1")

	persister := persister.NewLocal("/home/despondency/Downloads")

	v1Handlers.Post("/upload", v1.NewUploadFileHandler(persister).GetUploadFileHandler())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:8081"))
	if err != nil {
		log.Fatalf("cannot connect to mongo %v", err)
	}
	_ = client

	scanner.NewGoogleScanner()

	return app.Listen(fmt.Sprintf(":%d", port))
}

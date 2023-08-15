package application

import (
	"context"
	"fmt"
	"github.com/despondency/toggl-task/internal/handler/v1"
	"github.com/despondency/toggl-task/internal/persister"
	"github.com/despondency/toggl-task/internal/scanner"
	"github.com/despondency/toggl-task/internal/service"
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOpts := options.Client().SetHosts(
		[]string{"localhost:8081"},
	).SetAuth(
		options.Credential{
			Username: "mongouser",
			Password: "mongopass",
		},
	)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatalf("cannot connect to mongo %v", err)
	}
	mongoPersister := persister.NewMongoPersister(client)
	sc := scanner.NewGoogleScanner(ctx)
	p := persister.NewLocal("/home/despondency/Downloads")
	uploadSvc := service.NewMultiServicer(p, mongoPersister, sc)

	v1Handlers := app.Group("/v1")

	v1Handlers.Post("/receipt", v1.NewUploadReceiptHandler(uploadSvc).GetUploadFileHandler())
	v1Handlers.Get("/receipt", v1.NewGetReceiptResultHandler(uploadSvc).GetReceiptHandler())

	return app.Listen(fmt.Sprintf(":%d", port))
}

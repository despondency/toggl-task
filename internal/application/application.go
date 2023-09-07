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
	"os"
	"time"
)

type Application struct {
	port   int
	app    *fiber.App
	config *Config
}

type Builder struct {
	port   int
	config *Config
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) WithConfig(config *Config) *Builder {
	b.config = config
	return b
}

func (b *Builder) WithPort(port int) *Builder {
	b.port = port
	return b
}

func (b *Builder) Build() (*Application, error) {
	if b.port == 0 {
		return nil, fmt.Errorf("application port not set")
	}
	return &Application{
		port:   b.port,
		config: b.config,
	}, nil
}

func (a *Application) StartServer() error {
	app := fiber.New(
		fiber.Config{BodyLimit: 4 * 1024 * 1024}, //
	)
	a.app = app
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOpts := options.Client().SetHosts(
		[]string{a.config.DatabaseURI},
	).SetAuth(
		options.Credential{
			Username: a.config.DatabaseUser,
			Password: a.config.DatabasePassword,
		},
	)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatalf("cannot connect to mongo %v", err)
	}
	mongoPersister := persister.NewMongoPersister(client)
	sc := scanner.NewGoogleScanner(ctx)
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		log.Fatal(err)
	}
	defer func(path string) {
		err = os.RemoveAll(path)
		if err != nil {

		}
	}(dir)
	p := persister.NewLocal(dir)
	uploadSvc := service.NewMultiServicer(p, mongoPersister, sc)

	v1Handlers := app.Group("/v1")

	v1Handlers.Post("/receipt", v1.NewUploadReceiptHandler(uploadSvc).GetUploadFileHandler())
	v1Handlers.Get("/receipts-by-tags", v1.NewGetReceiptsByTagResultHandler(uploadSvc).GetReceiptsByTagHandler())
	v1Handlers.Get("/receipt", v1.NewGetReceiptResultHandler(uploadSvc).GetReceiptHandler())

	return app.Listen(fmt.Sprintf(":%d", a.port))
}

func (a *Application) StopServer() error {
	return a.app.Shutdown()
}

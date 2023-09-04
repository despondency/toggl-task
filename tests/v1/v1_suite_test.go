package v1

import (
	"github.com/despondency/toggl-task/internal/application"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	cfg := &application.Config{}
	cfg.DatabaseURI = "localhost:8081"
	cfg.Env = "test"
	cfg.DatabaseUser = "mongouser"
	cfg.DatabasePassword = "mongopass"
	app, err := application.NewBuilder().WithConfig(cfg).WithPort(8084).Build()
	go func() {
		err := app.StartServer()
		if err != nil {
			panic(err)
		}
	}()
	exitVal := m.Run()
	err = app.StopServer()
	if err != nil {
		panic(err)
	}
	os.Exit(exitVal)
}

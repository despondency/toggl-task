package v1

import (
	"github.com/despondency/toggl-task/internal/application"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	app := application.NewApplication(8084)
	go func() {
		err := app.StartServer()
		if err != nil {
			panic(err)
		}
	}()
	exitVal := m.Run()
	err := app.StopServer()
	if err != nil {
		panic(err)
	}
	os.Exit(exitVal)
}

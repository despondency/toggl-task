package main

import (
	"github.com/despondency/toggl-task/internal/application"
	"github.com/spf13/viper"
	"log"
)

func main() {
	cfg := &application.Config{}

	env := viper.GetString("ENVIRONMENT")
	if env == "" {
		env = "TEST"
	}
	dbUser := viper.GetString("DB_USER")
	if dbUser == "" {
		panic("db user is not set")
	}
	dbPassword := viper.GetString("DB_PASSWORD")
	if dbPassword == "" {
		panic("db password is not set")
	}
	dbUri := viper.GetString("DB_URI")
	if dbUri == "" {
		panic("db uri is not set")
	}
	cfg.Env = env
	cfg.DbUser = dbUser
	cfg.DbPassword = dbPassword
	cfg.DbURI = dbUri

	app, err := application.NewBuilder().WithConfig(cfg).WithPort(8080).Build()
	if err != nil {
		panic(err)
	}
	log.Fatal(app.StartServer())
}

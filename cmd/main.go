package main

import (
	"github.com/despondency/toggl-task/internal/application"
	"github.com/spf13/viper"
	"log"
)

func main() {
	cfg := &application.Config{}
	viper.AutomaticEnv()
	env := viper.GetString("ENVIRONMENT")
	if env == "" {
		env = "TEST"
	}
	cfg.Env = env
	dbUser := viper.GetString("DATABASE_USER")
	if dbUser == "" {
		panic("db user is not set")
	}
	cfg.DatabaseUser = dbUser
	dbPassword := viper.GetString("DATABASE_PASSWORD")
	if dbPassword == "" {
		panic("db password is not set")
	}
	cfg.DatabasePassword = dbPassword
	dbUri := viper.GetString("DATABASE_URI")
	if dbUri == "" {
		panic("db uri is not set")
	}
	cfg.DatabaseURI = dbUri
	port := viper.GetInt("PORT")
	if port != 0 {
		cfg.Port = port
	}
	app, err := application.NewBuilder().WithConfig(cfg).WithPort(port).Build()
	if err != nil {
		panic(err)
	}
	log.Fatal(app.StartServer())
}

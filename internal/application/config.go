package application

type Config struct {
	Env              string `mapstructure:"ENVIRONMENT"`
	DatabaseURI      string `mapstructure:"DATABASE_URI"`
	DatabasePassword string `mapstructure:"DATABASE_PASSWORD"`
	DatabaseUser     string `mapstructure:"DATABASE_USER"`
}

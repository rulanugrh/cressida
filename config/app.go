package config

import (
	"os"

	"github.com/joho/godotenv"
)

type App struct {
	Database struct {
		Name string
		Host string
		Port string
		User string
		Pass string
	}

	Server struct {
		Host   string
		Port   string
		Secret string
	}

	Opentelemetry struct {
		Host string
		Port string
	}
}

var app *App

func GetConfig() *App {
	if app == nil {
		app = initConfig()
	}

	return app
}
func initConfig() *App {
	conf := App{}
	if err := godotenv.Load(); err != nil {
		conf.Database.Host = ""
		conf.Database.Pass = ""
		conf.Database.Name = ""
		conf.Database.User = ""
		conf.Database.Port = ""

		conf.Server.Host = ""
		conf.Server.Port = ""
		conf.Server.Secret = ""
		return &conf
	}

	conf.Database.Host = os.Getenv("POSTGRESQL_HOST")
	conf.Database.Port = os.Getenv("POSTGRESQL_PORT")
	conf.Database.Name = os.Getenv("POSTGRESQL_NAME")
	conf.Database.User = os.Getenv("POSTGRESQL_USER")
	conf.Database.Pass = os.Getenv("POSTGRESQL_PASS")

	conf.Server.Secret = os.Getenv("SERVER_SECRET")
	conf.Server.Host = os.Getenv("SERVER_HOST")
	conf.Server.Port = os.Getenv("SERVER_PORT")

	return &conf
}

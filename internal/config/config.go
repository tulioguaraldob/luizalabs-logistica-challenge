package config

import (
	"os"

	"github.com/joho/godotenv"
)

// LoadEnvs loads all the variables listed in the environment variables file
func LoadEnvs() error {
	if Env == nil {
		if err := godotenv.Load(); err != nil {
			return err
		}

		env := &envVar{
			PostgresUser:     os.Getenv("POSTGRES_USER"),
			PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
			PostgresDb:       os.Getenv("POSTGRES_DB"),
			PostgresHost:     os.Getenv("POSTGRES_HOST"),
			PostgresPort:     os.Getenv("POSTGRES_PORT"),
			Port:             os.Getenv("PORT"),
		}

		Env = env
	}

	return nil
}

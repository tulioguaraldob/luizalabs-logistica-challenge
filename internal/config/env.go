package config

var Env *envVar

type envVar struct {
	PostgresUser     string
	PostgresPassword string
	PostgresDb       string
	PostgresHost     string
	PostgresPort     string
	Port             string
}

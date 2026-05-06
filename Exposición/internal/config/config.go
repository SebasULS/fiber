package config

import "os"

// Config contiene todas las variables de configuración de la app
type Config struct {
	Port string
	Env  string
}

// Load lee variables de entorno y devuelve una Config con valores por defecto
func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":3000"
	}

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	return &Config{
		Port: port,
		Env:  env,
	}
}

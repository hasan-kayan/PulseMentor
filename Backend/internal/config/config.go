package config

import "time"

type Config struct {
	Env string

	ServerHost string
	ServerPort string

	DatabaseURL string

	JWTSecret        string
	JWTIssuer        string
	AccessTokenTTL   time.Duration
	RefreshTokenTTL  time.Duration
	BcryptCost       int
	AllowInsecureDev bool
}

func (c Config) ServerAddr() string {
	host := c.ServerHost
	if host == "" {
		host = "0.0.0.0"
	}
	port := c.ServerPort
	if port == "" {
		port = "8080"
	}
	return host + ":" + port
}

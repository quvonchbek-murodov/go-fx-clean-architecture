package config

import (
	"fmt"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

var (
	cfg  *Config
	once sync.Once
)

type Config struct {
	AppName  string `env:"APP_NAME" env-default:"golang-project-structure"`
	AppEnv   string `env:"APP_ENV" env-default:"development"`
	AppDebug bool   `env:"APP_DEBUG" env-default:"true"`

	HTTPPort int `env:"HTTP_PORT" env-default:"8080"`
	GRPCPort int `env:"GRPC_PORT" env-default:"9090"`

	DB DB `env-prefix:"DB_"`
}

type DB struct {
	Driver       string   `env:"DRIVER" env-default:"postgres"`
	Host         string   `env:"HOST" env-default:"localhost"`
	Port         string   `env:"PORT" env-default:"5432"`
	User         string   `env:"USER" env-default:"postgres"`
	Password     string   `env:"PASSWORD" env-default:"postgres"`
	Name         string   `env:"NAME" env-default:"app"`
	SSLMode      string   `env:"SSLMODE" env-default:"disable"`
	ReplicaHosts []string `env:"REPLICA_HOSTS"`
}

func (d DB) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode)
}

func (d DB) ReplicaDSNs() []string {
	dsns := make([]string, 0, len(d.ReplicaHosts))
	for _, host := range d.ReplicaHosts {
		dsns = append(dsns, fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			host, d.Port, d.User, d.Password, d.Name, d.SSLMode))
	}
	return dsns
}

func LoadConfig() *Config {
	once.Do(func() {
		_ = godotenv.Load()

		cfg = new(Config)
		if err := cleanenv.ReadEnv(cfg); err != nil {
			panic(fmt.Errorf("failed to read config: %w", err))
		}
	})
	return cfg
}

package config

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type DataBaseConfig struct {
	Host           string `yaml:"host" default:"localhost"`
	Port           uint16 `yaml:"port" default:"5432"`
	User           string `yaml:"user" default:"postgres"`
	Password       string `yaml:"password"`
	DBName         string `yaml:"db_name" default:"tictactoe_db"`
	MaxConnections int    `yaml:"max_connections" default:"100"`
}

func NewDBConnection(cfg *DataBaseConfig) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s pool_max_conns=%d", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.MaxConnections)
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}

func Load() (*DataBaseConfig, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml"
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg DataBaseConfig
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

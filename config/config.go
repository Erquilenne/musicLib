package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Database DatabaseConfig `mapstructure:"database"`
	Logger   Logger         `mapstructure:"logger"`
	Server   ServerConfig   `mapstructure:"server"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}

type ServerConfig struct {
	Mode         string        `mapstructure:"mode"`
	AppVersion   string        `mapstructure:"app_version"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	Port         string        `mapstructure:"port"`
	PprofPort    string        `mapstructure:"pprof_port"`
	Debug        bool          `mapstructure:"debug"`
	ExternalMusicAPI string    `mapstructure:"external_music_api"`
}

type Logger struct {
	Development       bool   `mapstructure:"development"`
	DisableCaller     bool   `mapstructure:"disable_caller"`
	DisableStacktrace bool   `mapstructure:"disable_stacktrace"`
	Encoding          string `mapstructure:"encoding"`
	Level             string `mapstructure:"level"`
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func LoadConfig(filename string) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigName(filename)
	v.SetConfigType("json")
	v.AddConfigPath(".")

	v.AddConfigPath("./config")

	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, err
		}
		return nil, err
	}

	return v, nil
}

func ParseConfig(v *viper.Viper) (*Config, error) {
	var cfg Config
	err := v.Unmarshal(&cfg)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	return &cfg, nil
}

func GetConfigPath(configPath string) string {
	return "config"
}

func (p *PostgresConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		p.Host, p.Port, p.User, p.Password, p.DBName,
	)
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func LoadDatabaseConfig() (*PostgresConfig, error) {
	cfg := &PostgresConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "musiclib"),
	}

	return cfg, nil
}

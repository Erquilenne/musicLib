package config

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Database DatabaseConfig `mapstructure:"database"`
	Logger   Logger         `mapstructure:"logger"`
	Server   ServerConfig   `mapstructure:"server"`
	MusicApi MusicApiConfig `mapstructure:"music_api"`
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
}

type Logger struct {
	Development       bool   `mapstructure:"development"`
	DisableCaller     bool   `mapstructure:"disable_caller"`
	DisableStacktrace bool   `mapstructure:"disable_stacktrace"`
	Encoding          string `mapstructure:"encoding"`
	Level             string `mapstructure:"level"`
	EnableDebug       bool   `mapstructure:"enable_debug"`
}

type MusicApiConfig struct {
	URL string `mapstructure:"url"`
}

func LoadConfig(filename string) (*viper.Viper, error) {
	v := viper.New()

	// Загружаем .env файл
	v.SetConfigFile(".env")
	if err := v.ReadInConfig(); err != nil {
		log.Printf("Warning: .env file not found or error reading it: %v", err)
	}

	// Загружаем основной конфиг
	v.SetConfigName(filename)
	v.SetConfigType("json")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")

	// Автоматически читаем переменные окружения
	v.AutomaticEnv()
	v.SetEnvPrefix("MUSIC") // Префикс для переменных окружения
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // Заменяем точки на подчеркивания в именах

	// Устанавливаем соответствие между переменными окружения и конфигом
	v.BindEnv("music_api.url", "MUSIC_API_URL")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
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

func (p *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		p.Host, p.Port, p.User, p.Password, p.DBName,
	)
}

package config

import (
    "github.com/spf13/viper"
    "strings"
    "errors"
)

type Config struct {
    App      AppConfig
    Server   ServerConfig
    JWT      JWTConfig
    Database DatabaseConfig
}

type AppConfig struct {
    Name    string
    Version string
}

type ServerConfig struct {
    Port string
    Host string
}

type JWTConfig struct {
    SecretKey string `mapstructure:"secret_key"`
    ExpiresIn string `mapstructure:"expires_in"`
}

type DatabaseConfig struct {
    Host     string
    Port     string
    Name     string
    User     string
    Password string
}

func LoadConfig(path string) (*Config, error) {
    viper.AddConfigPath(path)
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")

    // Добавляем поддержку переменных окружения
    viper.AutomaticEnv()
    viper.SetEnvPrefix("APP") // Опционально: префикс для переменных окружения
    viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

    if err := viper.ReadInConfig(); err != nil {
        return nil, err
    }

    // Явно указываем, что JWT_SECRET должен быть взят из переменных окружения
    if err := viper.BindEnv("jwt.secret_key", "JWT_SECRET"); err != nil {
        return nil, err
    }

    var config Config
    if err := viper.Unmarshal(&config); err != nil {
        return nil, err
    }

    // Проверка наличия обязательных значений
    if config.JWT.SecretKey == "" {
        return nil, errors.New("JWT secret key is required")
    }

    return &config, nil
}

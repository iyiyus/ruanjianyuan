package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	DB        DBConfig
	JWT       JWTConfig
	App       AppConfig
	AdminPath string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type JWTConfig struct {
	Secret string
}

type AppConfig struct {
	Port string
}

func (d DBConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		d.User, d.Password, d.Host, d.Port, d.DBName)
}

func Load() *Config {
	// 先从 .env 文件读取
	envMap := loadEnvFile(".env")

	cfg := &Config{
		DB: DBConfig{
			Host:     getVal(envMap, "DB_HOST", "DB_HOST", "127.0.0.1"),
			Port:     getVal(envMap, "DB_PORT", "DB_PORT", "3306"),
			User:     getVal(envMap, "DB_USER", "DB_USER", "root"),
			Password: getVal(envMap, "DB_PASSWORD", "DB_PASSWORD", "yuchen"),
			DBName:   getVal(envMap, "DB_NAME", "DB_NAME", "ruanjianyuan"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "ruanjianyuan_secret_2024"),
		},
		App: AppConfig{
			Port: getVal(envMap, "APP_PORT", "APP_PORT", "1117"),
		},
		AdminPath: getVal(envMap, "ADMIN_PATH", "ADMIN_PATH", "admin"),
	}
	fmt.Printf("[Config] DB=%s@%s:%s/%s Port=%s\n", cfg.DB.User, cfg.DB.Host, cfg.DB.Port, cfg.DB.DBName, cfg.App.Port)
	return cfg
}

func loadEnvFile(path string) map[string]string {
	m := make(map[string]string)
	f, err := os.Open(path)
	if err != nil {
		return m
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "[") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			m[key] = val
		}
	}
	return m
}

func getVal(envMap map[string]string, key, envKey, def string) string {
	if v, ok := envMap[key]; ok && v != "" {
		return v
	}
	if envKey != "" {
		if v := os.Getenv(envKey); v != "" {
			return v
		}
	}
	return def
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

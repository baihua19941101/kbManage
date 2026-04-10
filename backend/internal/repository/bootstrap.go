package repository

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	HTTPAddr             string
	MySQLHost            string
	MySQLPort            int
	MySQLUser            string
	MySQLPassword        string
	MySQLDatabase        string
	MySQLParseTime       bool
	RedisAddr            string
	RedisPassword        string
	RedisDB              int
	JWTSecret            string
	AccessTokenTTL       time.Duration
	RefreshTokenTTL      time.Duration
	AdminSeedEnabled     bool
	AdminSeedUsername    string
	AdminSeedPassword    string
	AdminSeedDisplayName string
	AdminSeedEmail       string
	CORSAllowOrigins     []string
	CORSAllowMethods     []string
	CORSAllowHeaders     []string
	CORSExposeHeaders    []string
	CORSAllowCredentials bool
	CORSMaxAgeSeconds    int
}

type fileConfig struct {
	Server struct {
		HTTPAddr string `yaml:"http_addr"`
	} `yaml:"server"`
	MySQL struct {
		Host      string `yaml:"host"`
		Port      int    `yaml:"port"`
		User      string `yaml:"user"`
		Password  string `yaml:"password"`
		Database  string `yaml:"database"`
		ParseTime *bool  `yaml:"parse_time"`
	} `yaml:"mysql"`
	Redis struct {
		Addr     string `yaml:"addr"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
	} `yaml:"redis"`
	Security struct {
		JWTSecret       string `yaml:"jwt_secret"`
		AccessTokenTTL  string `yaml:"access_token_ttl"`
		RefreshTokenTTL string `yaml:"refresh_token_ttl"`
	} `yaml:"security"`
	Bootstrap struct {
		Admin struct {
			Enabled     *bool  `yaml:"enabled"`
			Username    string `yaml:"username"`
			Password    string `yaml:"password"`
			DisplayName string `yaml:"display_name"`
			Email       string `yaml:"email"`
		} `yaml:"admin"`
	} `yaml:"bootstrap"`
	CORS struct {
		AllowOrigins     []string `yaml:"allow_origins"`
		AllowMethods     []string `yaml:"allow_methods"`
		AllowHeaders     []string `yaml:"allow_headers"`
		ExposeHeaders    []string `yaml:"expose_headers"`
		AllowCredentials *bool    `yaml:"allow_credentials"`
		MaxAgeSeconds    int      `yaml:"max_age_seconds"`
	} `yaml:"cors"`
}

func LoadConfig() (Config, error) {
	configPath := getenv("CONFIG_FILE", "config/config.dev.yaml")
	candidates := []string{configPath}
	if configPath == "config/config.dev.yaml" {
		candidates = append(candidates, "backend/config/config.dev.yaml")
	}

	var (
		raw      []byte
		err      error
		usedPath string
	)
	for _, p := range candidates {
		raw, err = os.ReadFile(p)
		if err == nil {
			usedPath = p
			break
		}
	}
	if err != nil {
		return Config{}, fmt.Errorf("read config file failed (candidates: %v): %w", candidates, err)
	}

	var fc fileConfig
	if err := yaml.Unmarshal(raw, &fc); err != nil {
		return Config{}, fmt.Errorf("parse config file failed (%s): %w", usedPath, err)
	}

	accessTTL, err := time.ParseDuration(defaultIfEmpty(fc.Security.AccessTokenTTL, "15m"))
	if err != nil {
		return Config{}, fmt.Errorf("invalid security.access_token_ttl: %w", err)
	}

	refreshTTL, err := time.ParseDuration(defaultIfEmpty(fc.Security.RefreshTokenTTL, "168h"))
	if err != nil {
		return Config{}, fmt.Errorf("invalid security.refresh_token_ttl: %w", err)
	}

	cfg := Config{
		HTTPAddr:             defaultIfEmpty(fc.Server.HTTPAddr, ":8080"),
		MySQLHost:            defaultIfEmpty(fc.MySQL.Host, "127.0.0.1"),
		MySQLPort:            defaultIfZero(fc.MySQL.Port, 3306),
		MySQLUser:            defaultIfEmpty(fc.MySQL.User, "root"),
		MySQLPassword:        fc.MySQL.Password,
		MySQLDatabase:        defaultIfEmpty(fc.MySQL.Database, "kbmanage"),
		MySQLParseTime:       defaultBool(fc.MySQL.ParseTime, true),
		RedisAddr:            defaultIfEmpty(fc.Redis.Addr, "127.0.0.1:6379"),
		RedisPassword:        fc.Redis.Password,
		RedisDB:              fc.Redis.DB,
		JWTSecret:            defaultIfEmpty(fc.Security.JWTSecret, "dev-secret-change-me"),
		AccessTokenTTL:       accessTTL,
		RefreshTokenTTL:      refreshTTL,
		AdminSeedEnabled:     defaultBool(fc.Bootstrap.Admin.Enabled, true),
		AdminSeedUsername:    defaultIfEmpty(fc.Bootstrap.Admin.Username, "admin"),
		AdminSeedPassword:    defaultIfEmpty(fc.Bootstrap.Admin.Password, "Admin@123456"),
		AdminSeedDisplayName: defaultIfEmpty(fc.Bootstrap.Admin.DisplayName, "Administrator"),
		AdminSeedEmail:       defaultIfEmpty(fc.Bootstrap.Admin.Email, "admin@kbmanage.local"),
		CORSAllowOrigins:     defaultIfEmptySlice(fc.CORS.AllowOrigins, []string{"http://127.0.0.1:5000", "http://localhost:5000"}),
		CORSAllowMethods:     defaultIfEmptySlice(fc.CORS.AllowMethods, []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}),
		CORSAllowHeaders:     defaultIfEmptySlice(fc.CORS.AllowHeaders, []string{"Authorization", "Content-Type", "X-Request-Id"}),
		CORSExposeHeaders:    defaultIfEmptySlice(fc.CORS.ExposeHeaders, []string{"X-Request-Id"}),
		CORSAllowCredentials: defaultBool(fc.CORS.AllowCredentials, true),
		CORSMaxAgeSeconds:    defaultIfZero(fc.CORS.MaxAgeSeconds, 600),
	}

	// 支持通过环境变量覆盖关键项，便于 CI 或临时调试。
	cfg.HTTPAddr = getenv("HTTP_ADDR", cfg.HTTPAddr)
	cfg.MySQLHost = getenv("MYSQL_HOST", cfg.MySQLHost)
	cfg.MySQLPort = getenvInt("MYSQL_PORT", cfg.MySQLPort)
	cfg.MySQLUser = getenv("MYSQL_USER", cfg.MySQLUser)
	cfg.MySQLPassword = getenv("MYSQL_PASSWORD", cfg.MySQLPassword)
	cfg.MySQLDatabase = getenv("MYSQL_DATABASE", cfg.MySQLDatabase)
	cfg.RedisAddr = getenv("REDIS_ADDR", cfg.RedisAddr)
	cfg.RedisPassword = getenv("REDIS_PASSWORD", cfg.RedisPassword)
	cfg.RedisDB = getenvInt("REDIS_DB", cfg.RedisDB)
	cfg.JWTSecret = getenv("JWT_SECRET", cfg.JWTSecret)
	cfg.AdminSeedEnabled = getenvBool("ADMIN_SEED_ENABLED", cfg.AdminSeedEnabled)
	cfg.AdminSeedUsername = getenv("ADMIN_SEED_USERNAME", cfg.AdminSeedUsername)
	cfg.AdminSeedPassword = getenv("ADMIN_SEED_PASSWORD", cfg.AdminSeedPassword)
	cfg.AdminSeedDisplayName = getenv("ADMIN_SEED_DISPLAY_NAME", cfg.AdminSeedDisplayName)
	cfg.AdminSeedEmail = getenv("ADMIN_SEED_EMAIL", cfg.AdminSeedEmail)
	cfg.CORSAllowOrigins = getenvCSV("CORS_ALLOW_ORIGINS", cfg.CORSAllowOrigins)
	cfg.CORSAllowMethods = getenvCSV("CORS_ALLOW_METHODS", cfg.CORSAllowMethods)
	cfg.CORSAllowHeaders = getenvCSV("CORS_ALLOW_HEADERS", cfg.CORSAllowHeaders)
	cfg.CORSExposeHeaders = getenvCSV("CORS_EXPOSE_HEADERS", cfg.CORSExposeHeaders)
	cfg.CORSAllowCredentials = getenvBool("CORS_ALLOW_CREDENTIALS", cfg.CORSAllowCredentials)
	cfg.CORSMaxAgeSeconds = getenvInt("CORS_MAX_AGE_SECONDS", cfg.CORSMaxAgeSeconds)

	return cfg, nil
}

func (c Config) MySQLDSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=%t&loc=Local",
		c.MySQLUser,
		c.MySQLPassword,
		c.MySQLHost,
		c.MySQLPort,
		c.MySQLDatabase,
		c.MySQLParseTime,
	)
}

func NewGormDB(cfg Config) (*gorm.DB, error) {
	gcfg := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	}
	return gorm.Open(mysql.Open(cfg.MySQLDSN()), gcfg)
}

func getenv(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}

func getenvInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}

func getenvBool(key string, fallback bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return fallback
	}
	return b
}

func getenvCSV(key string, fallback []string) []string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback
	}
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		item := strings.TrimSpace(p)
		if item != "" {
			out = append(out, item)
		}
	}
	if len(out) == 0 {
		return fallback
	}
	return out
}

func defaultIfEmpty(v, fallback string) string {
	if v == "" {
		return fallback
	}
	return v
}

func defaultIfZero(v, fallback int) int {
	if v == 0 {
		return fallback
	}
	return v
}

func defaultBool(v *bool, fallback bool) bool {
	if v == nil {
		return fallback
	}
	return *v
}

func defaultIfEmptySlice(v []string, fallback []string) []string {
	if len(v) == 0 {
		return fallback
	}
	return v
}

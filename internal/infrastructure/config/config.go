package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App        AppConfig
	Server     ServerConfig
	Database   DatabaseConfig
	Redis      RedisConfig
	RabbitMQ   RabbitMQConfig
	JWT        JWTConfig
	CORS       CORSConfig
	RateLimit  RateLimitConfig
	Log        LogConfig
	Metrics    MetricsConfig
	Security   SecurityConfig
	Pagination PaginationConfig
}

type AppConfig struct {
	Name     string
	Env      string
	Port     int
	Debug    bool
	Timezone string
}

type ServerConfig struct {
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type DatabaseConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	Name            string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
	PoolSize int
}

type RabbitMQConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	VHost    string
}

type JWTConfig struct {
	Secret             string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
}

type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
	ExposedHeaders []string
	MaxAge         time.Duration
}

type RateLimitConfig struct {
	Enabled           bool
	RequestsPerSecond float64
	Burst             int
}

type LogConfig struct {
	Level  string
	Format string
	Output string
}

type MetricsConfig struct {
	Enabled bool
	Port    int
}

type SecurityConfig struct {
	BcryptCost        int
	PasswordMinLength int
}

type PaginationConfig struct {
	DefaultPageSize int
	MaxPageSize     int
}

func Load() (*Config, error) {
	v := viper.New()

	// Set config file
	v.SetConfigFile(".env")
	v.SetConfigType("env")

	// Auto env
	v.AutomaticEnv()

	// Read config
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	// Parse durations
	serverReadTimeout, _ := time.ParseDuration(v.GetString("SERVER_READ_TIMEOUT"))
	serverWriteTimeout, _ := time.ParseDuration(v.GetString("SERVER_WRITE_TIMEOUT"))
	serverIdleTimeout, _ := time.ParseDuration(v.GetString("SERVER_IDLE_TIMEOUT"))
	dbConnMaxLifetime, _ := time.ParseDuration(v.GetString("DB_CONN_MAX_LIFETIME"))
	jwtAccessExpiry, _ := time.ParseDuration(v.GetString("JWT_ACCESS_TOKEN_EXPIRY"))
	jwtRefreshExpiry, _ := time.ParseDuration(v.GetString("JWT_REFRESH_TOKEN_EXPIRY"))
	corsMaxAge, _ := time.ParseDuration(v.GetString("CORS_MAX_AGE"))

	config := &Config{
		App: AppConfig{
			Name:     v.GetString("APP_NAME"),
			Env:      v.GetString("APP_ENV"),
			Port:     v.GetInt("APP_PORT"),
			Debug:    v.GetBool("APP_DEBUG"),
			Timezone: v.GetString("APP_TIMEZONE"),
		},
		Server: ServerConfig{
			ReadTimeout:  serverReadTimeout,
			WriteTimeout: serverWriteTimeout,
			IdleTimeout:  serverIdleTimeout,
		},
		Database: DatabaseConfig{
			Host:            v.GetString("DB_HOST"),
			Port:            v.GetInt("DB_PORT"),
			User:            v.GetString("DB_USER"),
			Password:        v.GetString("DB_PASSWORD"),
			Name:            v.GetString("DB_NAME"),
			SSLMode:         v.GetString("DB_SSLMODE"),
			MaxOpenConns:    v.GetInt("DB_MAX_OPEN_CONNS"),
			MaxIdleConns:    v.GetInt("DB_MAX_IDLE_CONNS"),
			ConnMaxLifetime: dbConnMaxLifetime,
		},
		Redis: RedisConfig{
			Host:     v.GetString("REDIS_HOST"),
			Port:     v.GetInt("REDIS_PORT"),
			Password: v.GetString("REDIS_PASSWORD"),
			DB:       v.GetInt("REDIS_DB"),
			PoolSize: v.GetInt("REDIS_POOL_SIZE"),
		},
		RabbitMQ: RabbitMQConfig{
			Host:     v.GetString("RABBITMQ_HOST"),
			Port:     v.GetInt("RABBITMQ_PORT"),
			User:     v.GetString("RABBITMQ_USER"),
			Password: v.GetString("RABBITMQ_PASSWORD"),
			VHost:    v.GetString("RABBITMQ_VHOST"),
		},
		JWT: JWTConfig{
			Secret:             v.GetString("JWT_SECRET"),
			AccessTokenExpiry:  jwtAccessExpiry,
			RefreshTokenExpiry: jwtRefreshExpiry,
		},
		CORS: CORSConfig{
			AllowedOrigins: v.GetStringSlice("CORS_ALLOWED_ORIGINS"),
			AllowedMethods: v.GetStringSlice("CORS_ALLOWED_METHODS"),
			AllowedHeaders: v.GetStringSlice("CORS_ALLOWED_HEADERS"),
			ExposedHeaders: v.GetStringSlice("CORS_EXPOSED_HEADERS"),
			MaxAge:         corsMaxAge,
		},
		RateLimit: RateLimitConfig{
			Enabled:           v.GetBool("RATE_LIMIT_ENABLED"),
			RequestsPerSecond: v.GetFloat64("RATE_LIMIT_REQUESTS_PER_SECOND"),
			Burst:             v.GetInt("RATE_LIMIT_BURST"),
		},
		Log: LogConfig{
			Level:  v.GetString("LOG_LEVEL"),
			Format: v.GetString("LOG_FORMAT"),
			Output: v.GetString("LOG_OUTPUT"),
		},
		Metrics: MetricsConfig{
			Enabled: v.GetBool("METRICS_ENABLED"),
			Port:    v.GetInt("METRICS_PORT"),
		},
		Security: SecurityConfig{
			BcryptCost:        v.GetInt("BCRYPT_COST"),
			PasswordMinLength: v.GetInt("PASSWORD_MIN_LENGTH"),
		},
		Pagination: PaginationConfig{
			DefaultPageSize: v.GetInt("DEFAULT_PAGE_SIZE"),
			MaxPageSize:     v.GetInt("MAX_PAGE_SIZE"),
		},
	}

	return config, nil
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.Name,
		c.Database.SSLMode,
	)
}

func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port)
}

func (c *Config) GetRabbitMQURL() string {
	return fmt.Sprintf(
		"amqp://%s:%s@%s:%d%s",
		c.RabbitMQ.User,
		c.RabbitMQ.Password,
		c.RabbitMQ.Host,
		c.RabbitMQ.Port,
		c.RabbitMQ.VHost,
	)
}

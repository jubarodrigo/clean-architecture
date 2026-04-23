package configs

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Config espelha o padrão do módulo goexpert (Viper) com variáveis de ambiente.
type Config struct {
	DatabaseURL       string `mapstructure:"DATABASE_URL"`
	MigrationsDir     string `mapstructure:"MIGRATIONS_DIR"`
	WebServerPort     string `mapstructure:"WEB_SERVER_PORT"`
	GRPCServerPort    string `mapstructure:"GRPC_SERVER_PORT"`
	GraphQLServerPort string `mapstructure:"GRAPHQL_SERVER_PORT"`
}

// Load lê o arquivo .env (opcional) e sobrescreve com variáveis de ambiente.
func Load(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigName(".env")
	v.SetConfigType("env")
	if path != "" {
		v.AddConfigPath(path)
	}
	v.AddConfigPath(".")
	v.AddConfigPath("./configs")
	v.AddConfigPath("./cmd/ordersystem")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	_ = v.ReadInConfig()
	setDefaults(v)

	var c Config
	if err := v.Unmarshal(&c); err != nil {
		return nil, err
	}
	applyEnvOS(&c)
	applyStringDefaults(&c)
	if c.DatabaseURL == "" {
		return nil, fmt.Errorf("config: DATABASE_URL é obrigatório")
	}
	return &c, nil
}

func applyEnvOS(c *Config) {
	if s := os.Getenv("DATABASE_URL"); s != "" {
		c.DatabaseURL = s
	}
	if s := os.Getenv("MIGRATIONS_DIR"); s != "" {
		c.MigrationsDir = s
	}
	if s := os.Getenv("WEB_SERVER_PORT"); s != "" {
		c.WebServerPort = s
	}
	if s := os.Getenv("GRPC_SERVER_PORT"); s != "" {
		c.GRPCServerPort = s
	}
	if s := os.Getenv("GRAPHQL_SERVER_PORT"); s != "" {
		c.GraphQLServerPort = s
	}
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("DATABASE_URL", "postgres://orders:orders@localhost:5432/orders?sslmode=disable")
	v.SetDefault("MIGRATIONS_DIR", "migrations")
	v.SetDefault("WEB_SERVER_PORT", ":8080")
	v.SetDefault("GRPC_SERVER_PORT", "50051")
	v.SetDefault("GRAPHQL_SERVER_PORT", "8081")
}

func applyStringDefaults(c *Config) {
	if c.MigrationsDir == "" {
		c.MigrationsDir = "migrations"
	}
	if c.WebServerPort == "" {
		c.WebServerPort = ":8080"
	}
	if c.GRPCServerPort == "" {
		c.GRPCServerPort = "50051"
	}
	if c.GraphQLServerPort == "" {
		c.GraphQLServerPort = "8081"
	}
}

// GRPCAddress retorna endereço no formato ":porta" para net.Listen.
func (c *Config) GRPCAddress() string {
	if strings.HasPrefix(c.GRPCServerPort, ":") {
		return c.GRPCServerPort
	}
	return ":" + c.GRPCServerPort
}

// GraphQLAddress retorna endereço no formato ":porta" para http.ListenAndServe.
func (c *Config) GraphQLAddress() string {
	p := c.GraphQLServerPort
	if strings.HasPrefix(p, ":") {
		return p
	}
	return ":" + p
}

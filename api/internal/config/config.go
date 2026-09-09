package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Port               string
	SupabaseURL        string
	SupabaseAnonKey    string
	SupabaseServiceKey string
	InternalAPIKey     string
}

type yamlConfig struct {
	Supabase struct {
		URL        string `yaml:"url"`
		AnonKey    string `yaml:"anon_key"`
		ServiceKey string `yaml:"service_key"`
	} `yaml:"supabase"`

	Server struct {
		GoPort int `yaml:"go_port"`
	} `yaml:"server"`
}

func resolveEnv(value string) string {
	if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
		envName := value[2 : len(value)-1]
		return os.Getenv(envName)
	}
	return value
}

func Load() *Config {
	if err := godotenv.Load("../.env"); err != nil {
		log.Printf("WARNING: failed to load .env file: %v; using environment variables instead", err)
	}

	yamlPath := os.Getenv("CONFIG_YAML")
	if yamlPath == "" {
		yamlPath = "../config.yaml"
	}

	data, err := os.ReadFile(yamlPath)
	if err != nil {
		log.Fatalf(
			"CONFIG ERROR: failed to read config file %q: %v",
			yamlPath,
			err,
		)
	}

	var yc yamlConfig
	if err := yaml.Unmarshal(data, &yc); err != nil {
		log.Fatalf(
			"CONFIG ERROR: failed to parse config file %q: %v",
			yamlPath,
			err,
		)
	}

	port := os.Getenv("PORT")
	if port == "" && yc.Server.GoPort != 0 {
		port = fmt.Sprintf("%d", yc.Server.GoPort)
	}
	if port == "" {
		port = "8080"
	}

	cfg := &Config{
		Port:               port,
		SupabaseURL:        yc.Supabase.URL,
		SupabaseAnonKey:    resolveEnv(yc.Supabase.AnonKey),
		SupabaseServiceKey: resolveEnv(yc.Supabase.ServiceKey),
		InternalAPIKey:     os.Getenv("INTERNAL_API_KEY"),
	}

	if cfg.SupabaseURL == "" {
		log.Fatal(
			"CONFIG ERROR: Supabase URL is not configured. " +
				"Set 'supabase.url' in config.yaml",
		)
	}

	if cfg.SupabaseAnonKey == "" {
		log.Fatal(
			"CONFIG ERROR: Supabase anonymous key is not configured. " +
				"Set 'supabase.anon_key' in config.yaml",
		)
	}

	if cfg.SupabaseServiceKey == "" {
		log.Fatal(
			"CONFIG ERROR: Supabase service key is not configured. " +
				"Set 'supabase.service_key' in config.yaml",
		)
	}

	if cfg.InternalAPIKey == "" {
		log.Fatal(
			"CONFIG ERROR: INTERNAL_API_KEY is not configured. " +
				"Set INTERNAL_API_KEY in .env or environment variables",
		)
	}

	return cfg
}
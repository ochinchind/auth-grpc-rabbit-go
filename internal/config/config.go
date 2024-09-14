package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"time"
)

type Config struct {
	Env            string     `yaml:"env" env-default:"local"`
	StoragePath    string     `yaml:"storage_path" env-required:"true"`
	GRPC           GRPCConfig `yaml:"grpc"`
	MigrationsPath string
	TokenTTL       time.Duration `yaml:"token_ttl" env-default:"1h"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("config path is empty: " + err.Error())
	}

	return &cfg
}

// fetchConfigPath fetches config path from command line flag or environment variable.
// Priority: flag > env > default.
// Default value is empty string.
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}

// MustLoadPath loads configuration from a YAML file and panics if it fails.
func MustLoadPath(path string) *Config {
	// Open the file
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("failed to open config file: %v", err)
	}
	defer file.Close()

	// Create an instance of Config to hold the loaded configuration
	var cfg Config

	// Decode the YAML file into the Config struct
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		log.Fatalf("failed to decode config file: %v", err)
	}

	// Optionally log or validate the config here
	log.Printf("Config loaded from %s", path)

	return &cfg
}

package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"os"
	"time"
)

const (
	defaultGRPCPort = "443"
	authority       = "qrcode-generation-service"
	EnvLocal        = "local"
)

type (
	Config struct {
		Authority    string
		Environment  string
		GRPC         GRPCConfig         `mapstructure:"grpc"`
		Users        MicroserviceConfig `mapstructure:"userMicroservice"`
		Reservations MicroserviceConfig `mapstructure:"reservationMicroservice"`
	}

	GRPCConfig struct {
		Host    string        `mapstructure:"host"`
		Port    string        `mapstructure:"port"`
		Timeout time.Duration `mapstructure:"timeout"`
	}

	MicroserviceConfig struct {
		Host string `mapstructure:"host"`
		Port string `mapstructure:"port"`
	}
)

func Init(configsDir, envDir string) (*Config, error) {
	populateDefaults()
	loadEnvVariables(envDir)
	if err := parseConfigFile(configsDir); err != nil {
		return nil, err
	}

	var cfg Config
	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}

	setFromEnv(&cfg)

	return &cfg, nil
}

func unmarshal(cfg *Config) error {
	if err := viper.UnmarshalKey("userMicroservice", &cfg.Users); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("reservationMicroservice", &cfg.Reservations); err != nil {
		return err
	}
	return viper.UnmarshalKey("grpc", &cfg.GRPC)
}

func setFromEnv(cfg *Config) {

	cfg.GRPC.Host = os.Getenv("GRPC_HOST")

	cfg.Environment = EnvLocal
	cfg.Authority = authority
}

func parseConfigFile(folder string) error {
	viper.AddConfigPath(folder)
	viper.SetConfigName("main")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return viper.MergeInConfig()
}

func loadEnvVariables(envPath string) {
	err := godotenv.Load(envPath)

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

}

func populateDefaults() {
	viper.SetDefault("grpc.port", defaultGRPCPort)
}

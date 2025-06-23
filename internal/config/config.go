package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env" env-default:"local"`
	StoragePath string `yaml:"storage_path" env-requiared:"true"`
	HTTPServer  `yaml:"http_server"`
	PostgresDB  `yaml:"postgres_db"`
	Clients     ClientConfig `yaml:"clients"`
	AppSecret   string       `yaml:"app_secret"`
}

type HTTPServer struct {
	Address     string        `yaml:"address"`
	User        string        `yaml:"user" env-requiared:"true"`
	Password    string        `yaml:"password" env-requiared:"true" env:"HTTP_SERVER_PASSWORD"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type PostgresDB struct {
	User     string `yaml:"user"`
	Dbname   string `yaml:"dbname"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Sslmode  string `yaml:"sslmode"`
}

type Client struct {
	Address      string        `yaml:"addres" env-default:"localhost:44040"`
	Timeout      time.Duration `yaml:"timeout" env-default:"4s"`
	RetriesCount int           `yaml:"retriesCount" env-default:"2"`
	// Insecure     bool          `yaml:"insecure"`
}

type ClientConfig struct {
	SSO Client `yaml:"sso"`
}

func MustLoad() *Config {

	os.Setenv("CONFIG_PATH", "./config/local.yaml")

	// if err := godotenv.Load("local.env"); err != nil {
	// 	fmt.Printf("no .env file found, loading from system\n")
	// }

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	//check if file exists

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file %s does not exist", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}

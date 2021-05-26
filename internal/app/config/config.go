package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"time"
)

const (
	httpReadTimeout  = 5 * time.Second
	httpWriteTimeout = 10 * time.Second
	httpIdleTimeout  = time.Minute
	shutdownTimeout  = 5 * time.Second
)

type Config struct {
	Web struct {
		APIHost         string `yaml:"host"`
		APIPort         string `yaml:"port"`
		ReadTimeout     time.Duration
		WriteTimeout    time.Duration
		IdleTimeout     time.Duration
		ShutdownTimeout time.Duration
	} `yaml:"web"`

	Database struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Username string `yaml:"username"`
		DBName   string `yaml:"dbName"`
		Password string
	} `yaml:"postgres"`

	FileStorage struct {
		Endpoint        string `yaml:"endPoint"`
		BucketName      string `yaml:"bucketName"`
		AccessKey       string `yaml:"aws_access_key"`
		AwsRegion       string `yaml:"awsRegion"`
		SecretAccessKey string
	} `yaml:"fileStorage"`
}

func Init(configsDir string) (*Config, error) {
	var cfg Config
	content, err := ioutil.ReadFile(configsDir)
	if err != nil {
		return nil, err
	}

	if err = yaml.Unmarshal(content, &cfg); err != nil {
		return nil, err
	}

	setDefaultsAndEnv(&cfg)

	return &cfg, nil
}

func setDefaultsAndEnv(cfg *Config) {
	cfg.Web.ReadTimeout = httpReadTimeout
	cfg.Web.WriteTimeout = httpWriteTimeout
	cfg.Web.IdleTimeout = httpIdleTimeout
	cfg.Web.ShutdownTimeout = shutdownTimeout
	cfg.Database.Password = os.Getenv("POSTGRES_PASSWORD")
	cfg.FileStorage.SecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
}

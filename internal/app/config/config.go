package config

import (
	"expvar"
	"fmt"
	"github.com/ardanlabs/conf"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	_ "net/http/pprof"
	"os"
	"time"
)

var build = "develop"

type Config struct {
	conf.Version
	Web struct {
		APIHost         string        `conf:"default:127.0.0.1" yaml:"host"`
		APIPort         string        `conf:"default:4000" yaml:"port"`
		ReadTimeout     time.Duration `conf:"default:5s"`
		WriteTimeout    time.Duration `conf:"default:10s"`
		IdleTimeout     time.Duration `conf:"default:5s"`
		ShutdownTimeout time.Duration `conf:"default:5s"`
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

	SMTP struct {
		Host string `conf:"default:smtp.gmail.com" yaml:"host"`
		Port int    `conf:"default:587" yaml:"port"`
		From string `conf:"default:duman070601@gmail.com" yaml:"from"`
		Pass string
	}
}

func Init(cfg *Config, configsDir string) error {
	cfg.Version.SVN = build
	cfg.Version.Desc = "copyright information here"

	if err := conf.Parse(os.Args[1:], "BOOKING", cfg); err != nil {
		switch err {
		case conf.ErrHelpWanted:
			usage, err := conf.Usage("BOOKING", cfg)
			if err != nil {
				return errors.Wrap(err, "generating config usage")
			}
			fmt.Println(usage)
			return nil
		case conf.ErrVersionWanted:
			version, err := conf.VersionString("BOOKING", cfg)
			if err != nil {
				return errors.Wrap(err, "generating config version")
			}
			fmt.Println(version)
			return nil
		}
		return errors.Wrap(err, "parsing config")
	}

	setEnvVariables(cfg)

	content, err := ioutil.ReadFile(configsDir)
	if err != nil {
		return errors.Wrap(err, "reading file")
	}

	if err = yaml.Unmarshal(content, cfg); err != nil {
		return errors.Wrap(err, "unmarshalling config")
	}

	// App Starting
	expvar.NewString("build").Set(build)
	log.Printf("main: Started: Application initializing: version %q", build)

	out, err := conf.String(cfg)
	if err != nil {
		return errors.Wrap(err, "generating config for output")
	}
	log.Printf("main: Config:\n%v\n", out)

	return nil
}

func setEnvVariables(cfg *Config) {
	cfg.Database.Password = os.Getenv("POSTGRES_PASSWORD")
	cfg.FileStorage.SecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	cfg.SMTP.Pass = os.Getenv("GMAIL_PASSWORD")
}

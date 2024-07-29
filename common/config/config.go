package config

import (
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

// Add future services to config struct then define them below
type Config struct {
	AppPath string `yaml:"-"`
	AWS     AWS    `yaml:"aws"`
}

type AWS struct {
	AWSAccessKey       string `yaml:"aws_access_key"`
	AWSSecretAccessKey string `yaml:"aws_secret_access_key"`
	Session            string `yaml:"session"` //optional?
}

func NewConfig(appPath, envFile string) (Config, error) {
	var c Config
	if _, err := os.Stat(filepath.Join(appPath, envFile)); err != nil {
		return c, err
	}

	vp := viper.New()
	vp.AddConfigPath(appPath)
	vp.SetConfigType("yaml")
	vp.SetConfigName(envFile)
	if err := vp.ReadInConfig(); err != nil {
		return c, err
	}

	c.AppPath = appPath
	c.AWS.AWSAccessKey = vp.GetString("aws_access_key")
	c.AWS.AWSSecretAccessKey = vp.GetString("aws_secret_access_key")
	c.AWS.Session = vp.GetString("session")

	return c, nil
}

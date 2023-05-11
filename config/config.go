package config

import (
	"os"

	"github.com/spf13/viper"
)

var vp *viper.Viper

type App struct {
	PORT string `mapstructure:"PORT"`
}

type Config struct {
	APP App `mapstructure:"APP"`
}

var config Config

func ConfigInit() error {
	vp = viper.New()

	var curEnv string
	env := os.Getenv("GO_ENV")
	if env == "prod" {
		curEnv = env
	} else {
		curEnv = "build"
	}
	envFileName := curEnv + ".yaml"

	workingdir, err := os.Getwd()
	if err != nil {
		return err
	}

	configFilePath := workingdir + "/config/" + envFileName
	vp.SetConfigFile(configFilePath)

	err = vp.ReadInConfig()
	if err != nil {
		return err
	}

	err = vp.Unmarshal(&config)
	if err != nil {
		return err
	}

	return nil
}

func GetConfig() Config {
	return config
}

func GetViper() *viper.Viper {
	return vp
}

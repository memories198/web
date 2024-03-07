package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

var (
	MysqlConfig      = &Mysql{}
	RedisConfig      = &Redis{}
	LogFile          *os.File
	CookieExpireTime int
)

type Mysql struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Database string `yaml:"database"`
}
type Redis struct {
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Database int    `yaml:"database"`
	Timeout  int    `yaml:"timeout"`
}

func Init() error {
	configFile, err := os.ReadFile("config/conf.yaml")
	if err != nil {
		return err
	}

	Config := struct {
		*Mysql           `yaml:"mysql"`
		*Redis           `yaml:"redis"`
		File             string `yaml:"webLog"`
		CookieExpireTime *int   `yaml:"cookieExpireTime"`
	}{
		Mysql:            MysqlConfig,
		Redis:            RedisConfig,
		CookieExpireTime: &CookieExpireTime,
	}

	err = yaml.Unmarshal(configFile, &Config)
	if err != nil {
		return err
	}
	LogFile, err = os.OpenFile(Config.File, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	fmt.Println(CookieExpireTime)
	return nil
}

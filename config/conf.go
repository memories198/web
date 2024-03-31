package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

var (
	MysqlConfig      = &Mysql{}
	RedisConfig      = &Redis{}
	WebLogFile       *os.File
	GormLogFile      *os.File
	CookieExpireTime = 3600000
)

type Mysql struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Database string `yaml:"database"`
}
type Redis struct {
	Password string `yaml:"password" `
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
		WebLogFile       string `yaml:"webLog"`
		GormLogFile      string `yaml:"gormLogFile"`
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
	WebLogFile, err = os.OpenFile(Config.WebLogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	GormLogFile, err = os.OpenFile(Config.GormLogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	return nil
}

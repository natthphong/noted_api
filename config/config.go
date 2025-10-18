package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Env      string
	EnableS3 bool
	Server   Server
	Log      Log
	HTTP     HTTP
	DBConfig DBConfig
}

type S3Config struct {
	BucketName string
	Key        string
}

type Server struct {
	Name string
	Port string
}

type Log struct {
	Level string
}

type DBConfig struct {
	Host            string
	Port            string
	Username        string
	Password        string
	Name            string
	MaxOpenConn     int32
	MaxConnLifeTime int64
}

type HTTP struct {
	TimeOut            time.Duration
	MaxIdleConn        int
	MaxIdleConnPerHost int
	MaxConnPerHost     int
	CertFile           []byte
	KeyFile            []byte
}

func InitConfig() (*Config, error) {
	env := os.Getenv("ENV")
	if strings.EqualFold(env, "") {
		configPath, ok := os.LookupEnv("API_CONFIG_PATH")
		if !ok {
			configPath = "./config"
		}
		configName, ok := os.LookupEnv("API_CONFIG_NAME")
		if !ok {
			configName = "config"
		}
		viper.SetConfigName(configName)
		viper.AddConfigPath(configPath)
		if err := viper.ReadInConfig(); err != nil {
			fmt.Println("config file not found. using default/env config: " + err.Error())
		}
	}

	// defaults
	viper.SetDefault("ENV", env)
	viper.SetDefault("LOG.LEVEL", "debug")
	viper.SetDefault("DBCONFIG.MAXOPENCONN", "4")
	viper.SetDefault("DBCONFIG.MAXCONNLIFETIME", "300")
	viper.SetDefault("HTTP.TIMEOUT", "10s")
	viper.SetDefault("HTTP.MAXIDLECONN", 100)
	viper.SetDefault("HTTP.MAXIDLECONNPERHOST", 100)
	viper.SetDefault("HTTP.MAXCONNPERHOST", 100)

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		return nil, err
	}
	return &c, nil
}

func InitTimeZone() {
	ict, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		panic(err)
	}
	time.Local = ict
}

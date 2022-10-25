package conf

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

const (
	EnvironmentLocal = "LOCAL"
	EnvironmentDev   = "DEV"
)

// Config db environment
type Config struct {
	MySQL struct {
		Host     string
		Port     int64
		User     string
		Password string
		DB       string

		MaxIdleConns int
		MaxOpenConns int
	}

	Environment string
}

var EnvConfig *Config

func Get() *Config {
	return EnvConfig
}

// InitConfig init function
func InitConfig() error {
	EnvConfig = &Config{}

	var err error
	EnvConfig.Environment = os.Getenv("ENVIRONMENT")
	//log.Fatal("Environment: ", EnvConfig.Environment)
	if EnvConfig.Environment == EnvironmentLocal {
		err = godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
			return err
		}
	}
	EnvConfig.MySQL.Host = os.Getenv("MYSQL_HOST")
	mysqlPort, err := strconv.ParseInt(os.Getenv("MYSQL_PORT"), 10, 64)
	if err != nil {
		log.Fatal("Error when parse config MYSQL_PORT, detail: ", err)
		return err
	}
	EnvConfig.MySQL.Port = mysqlPort
	EnvConfig.MySQL.User = os.Getenv("MYSQL_USER")
	EnvConfig.MySQL.Password = os.Getenv("MYSQL_PASSWORD")
	EnvConfig.MySQL.DB = os.Getenv("MYSQL_DATABASE")
	if os.Getenv("MYSQL_MAX_IDLE_CONNS") != "" {
		maxIdleConns, err := strconv.ParseInt(os.Getenv("MYSQL_MAX_IDLE_CONNS"), 10, 64)
		if err != nil {
			log.Fatal("Error when parse config MYSQL_MAX_IDLE_CONNS, detail: ", err)
			return err
		}
		EnvConfig.MySQL.MaxIdleConns = int(maxIdleConns)
	}
	if os.Getenv("MYSQL_MAX_OPEN_CONNS") != "" {
		maxOpenConns, err := strconv.ParseInt(os.Getenv("MYSQL_MAX_OPEN_CONNS"), 10, 64)
		if err != nil {
			log.Fatal("Error when parse config MYSQL_MAX_OPEN_CONNS, detail: ", err)
			return err
		}
		EnvConfig.MySQL.MaxOpenConns = int(maxOpenConns)
	}
	return nil
}

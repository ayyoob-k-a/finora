package configs

import "os"

type Mail struct {
	SecretKey string `mapstructure:"SECRETKEY"`
	From      string `mapstructure:"FROM"`
	URL       string `mapstructure:"URL"`
}


type Config struct {
	DBNAME   string
	DBUSER   string
	PASSWORD string
	HOST     string
	PORT     string
	Mail    Mail
}

func GetConfig() Config {
	return Config{
		DBNAME: os.Getenv("DBNAME"),
		DBUSER: os.Getenv("DBUSER"),
		PASSWORD: os.Getenv("PASSWORD"),
		HOST: os.Getenv("HOST"),
		PORT: os.Getenv("PORT"),
		Mail: Mail{
			SecretKey: os.Getenv("SECRETKEY"),
			
		},
	}
}
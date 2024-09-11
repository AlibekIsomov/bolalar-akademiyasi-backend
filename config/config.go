package config

import (
	"fmt"
)

type Config struct {
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     string
}

func LoadConfig() *Config {
	return &Config{
		DBUser:     "postgres",
		DBPassword: "y@suk321",
		DBName:     "academy",
		DBHost:     "localhost",
		DBPort:     "5432",
	}
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		c.DBHost, c.DBUser, c.DBPassword, c.DBName, c.DBPort)
}

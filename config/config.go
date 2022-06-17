package config

import (
	"database/sql"
	"encoding/json"
	lr "go-inventory/util/logger"
	"os"
	"time"

	"github.com/joeshaw/envdecode"
	"github.com/joho/godotenv" // package used to read the .env file
	_ "github.com/lib/pq"      // postgres golang driver
)

type Conf struct {
	Debug  bool `env:"DEBUG,required"`
	Server serverConf
}

type serverConf struct {
	Port         int           `env:"SERVER_PORT,required"`
	TimeoutRead  time.Duration `env:"SERVER_TIMEOUT_READ,required"`
	TimeoutWrite time.Duration `env:"SERVER_TIMEOUT_WRITE,required"`
	TimeoutIdle  time.Duration `env:"SERVER_TIMEOUT_IDLE,required"`
}

// var logger Logger:= lr.New(false)

func AppConfig() *Conf {
	logger := lr.New(true)
	err := godotenv.Load()

	if err != nil {
		logger.Info().Msgf("Error loading .env file")
	}

	var c Conf
	if err := envdecode.StrictDecode(&c); err != nil {
		logger.Info().Msgf("Failed to decode: %s", err)
	}
	return &c
}

func CreateConnection() *sql.DB {
	appConf := AppConfig()
	logger := lr.New(appConf.Debug)
	logger.Info().Msgf("Connection %s\n", os.Getenv("POSTGRES_URL"))

	//Connection to db
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))

	if err != nil {
		panic(err)
	}

	// check the connection
	err = db.Ping()

	if err != nil {
		panic(err)
	}

	logger.Info().Msgf("Success db connection!")
	return db
}

type NullString struct {
	sql.NullString
}

func (s NullString) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(s.String)
}

func (s *NullString) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		s.String, s.Valid = "", false
		return nil
	}
	s.String, s.Valid = string(data), true
	return nil
}

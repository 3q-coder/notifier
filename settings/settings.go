package settings

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

var (
	// db
	DBDialect string
	DBHost    string
	DBName    string
	DBPort    string
	DBUser    string
	DBPw      string
	// web
	APIKey     string
	WriteWait  time.Duration
	PongWait   time.Duration
	PingPeriod time.Duration
)

func Init(dotenvFileName string) error {
	if dotenvFileName != "" {
		err := godotenv.Load(dotenvFileName)
		if err != nil {
			return err
		}
	}

	// data base settings
	DBDialect = getenv("DB_DIALECT", "sqlite")
	DBHost = os.Getenv("DB_HOST")
	DBName = getenv("DB_NAME", "local.db")
	DBPort = os.Getenv("DB_PORT")
	DBUser = os.Getenv("DB_USER")
	DBPw = os.Getenv("DB_PW")

	// web settings
	APIKey = os.Getenv("API_KEY")

	t, err := strconv.Atoi(os.Getenv("WRITE_WAIT"))
	if err != nil {
		return err
	}
	WriteWait = time.Duration(t) * time.Second
	t, err = strconv.Atoi(os.Getenv("PONG_WAIT"))
	if err != nil {
		return err
	}
	PongWait = time.Duration(t) * time.Second
	t, err = strconv.Atoi(os.Getenv("PING_PERIOD"))
	if err != nil {
		return err
	}
	PingPeriod = time.Duration(t) * time.Second

	return nil
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

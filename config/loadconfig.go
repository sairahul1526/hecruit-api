package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// LoadConfig - load .env file from given path for local, else will be getting from env var
func LoadConfig() {
	// load .env file from given path for local, else will be getting from env var
	if !strings.EqualFold(os.Getenv("prod"), "true") {
		err := godotenv.Load(".test-env")
		if err != nil {
			panic("Error loading .env file")
		}
	}

	DBConfig = os.Getenv("DB_CONFIG")
	DBConnectionPool, _ = strconv.Atoi("DB_CONNECTION_POOL")
	Log, _ = strconv.ParseBool(os.Getenv("LOG"))
	Migrate, _ = strconv.ParseBool(os.Getenv("MIGRATE"))
	JWTSecret = []byte(os.Getenv("JWT_SECRET"))

	// s3
	S3Endpoint = os.Getenv("S3_ENDPOINT")
	S3ID = os.Getenv("S3_ID")
	S3Secret = os.Getenv("S3_SECRET")
	S3Region = os.Getenv("S3_REGION")
	S3Bucket = os.Getenv("S3_BUCKET")
	S3MediaURL = os.Getenv("S3_MEDIA_URL")

	// aws
	AWSAccessKey = os.Getenv("AMAZON_ACCESS_KEY")
	AWSSecretKey = os.Getenv("AMAZON_SECRET_KEY")
}

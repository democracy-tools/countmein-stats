package env

import (
	"os"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func Load() {

	err := godotenv.Load("internal/env/.env")
	if err != nil {
		log.Fatalf("failed to load environment file with %q", err)
	}
}

func GetEnvSensitive(variable string) string {

	res := os.Getenv(variable)
	if res != "" {
		log.Debugf("%q: [sensitive]", variable)
	}

	return res
}

func GetProject() string {

	return getEnvOrExit("GCP_PROJECT_ID")
}

func GetBQDataset() string {

	return getEnvOrExit("BQ_DATASET")
}

func GetBucket() string {

	return getEnvOrExit("STORAGE_BUCKET")
}

func getEnvOrExit(variable string) string {

	res := os.Getenv(variable)
	if res == "" {
		log.Fatalf("Please, set %q", variable)
	}
	log.Debugf("%q: %q", variable, res)

	return res
}

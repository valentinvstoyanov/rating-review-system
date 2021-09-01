package env

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

const (
	envVarName = "ENV"
	envDirPath = "./env/"
	envFileExt = ".env"
	prodEnv    = "prod"
	devEnv     = "dev"
)

func GetEnvVar(name string) string {
	env := getEnv()
	envFileName := fmt.Sprintf("%s%s%s", envDirPath, env, envFileExt)

	if err := godotenv.Load(envFileName); err != nil {
		log.Fatalf("Failed to load %s.env file: %s", env, err)
	}

	return os.Getenv(name)
}

func IsProd() bool {
	return matchEnv(prodEnv)
}

func IsDev() bool {
	return matchEnv(devEnv)
}

func matchEnv(env string) bool {
	return getEnv() == env
}

func getEnv() string {
	env := os.Getenv(envVarName)
	if len(env) == 0 {
		panic(fmt.Sprintf("Missing %s environment variable", envVarName))
	}
	return env
}

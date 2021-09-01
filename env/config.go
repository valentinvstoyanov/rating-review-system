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

	res, ok := os.LookupEnv(name)
	if !ok {
		missingEnvVar(name)
	}

	return res
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
		missingEnvVar(envVarName)
	}
	return env
}

func missingEnvVar(name string) {
	panic(fmt.Sprintf("Missing %s environment variable", name))
}

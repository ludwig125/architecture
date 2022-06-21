package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	if err := runActorAPI(); err != nil {
		log.Panicf("failed to runActorAPI: %v", err)
	}
}

func runActorAPI() error {
	var repository ActorRepository
	var err error
	dbType := os.Getenv("DB_TYPE")
	switch dbType {
	case "sqlite":
		dbName := mustGetenv("DB_NAME")
		log.Println("use sqlite. database name:", dbName)
		repository, err = NewSQLiteActorRepository(dbName)
		if err != nil {
			return fmt.Errorf("failed to NewMySQLActorRepository: %v", err)
		}
	// case "mysql":
	// 	log.Println("use mysql")
	// 	repository, err = NewMySQLActorRepository("arch_db")
	// 	if err != nil {
	// 		log.Panicf("failed to NewMySQLActorRepository: %v", err)
	// 	}
	default:
		return fmt.Errorf("invalid dbType: %s", dbType)
	}

	exRepository, err := NewExcludeRepository(useEnvOrDefault("EXCLUDE_ACTOR_FILE", "exclude_actors.txt"))
	if err != nil {
		return fmt.Errorf("failed to NewExcludeRepository: %v", err)
	}

	config := Config{Port: useEnvOrDefault("SERVER_PORT", "8080")}

	service := NewActorService(config, repository, exRepository)

	server := NewServer(config, service)
	return server.Run()
}

type Config struct {
	Port string
}

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Panicf("environment variable '%s' not set", k)
	}

	return v
}

func useEnvOrDefault(key, def string) string {
	v := def
	if fromEnv := os.Getenv(key); fromEnv != "" {
		v = fromEnv
	}
	log.Printf("%s environment variable set", key)
	return v
}

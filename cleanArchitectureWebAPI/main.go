package main

import (
	"log"
	"os"
)

func main() {
	var repository ActorRepository
	var err error
	dbType := os.Getenv("DB_TYPE")
	switch dbType {
	case "sqlite":
		log.Println("use sqlite")
		repository, err = NewSQLiteActorRepository("arch_db")
		if err != nil {
			log.Panicf("failed to NewMySQLActorRepository: %v", err)
		}
	// case "mysql":
	// 	log.Println("use mysql")
	// 	repository, err = NewMySQLActorRepository("arch_db")
	// 	if err != nil {
	// 		log.Panicf("failed to NewMySQLActorRepository: %v", err)
	// 	}
	default:
		log.Panicf("invalid dbType: %s", dbType)
	}

	// p := Actor{
	// 	ID:   4,
	// 	Name: "Smith",
	// 	Age:  43,
	// }
	// if err := repository.Update(p); err != nil {
	// 	log.Panic(err)
	// }
	// if err := repository.DeleteByID(4); err != nil {
	// 	log.Panic(err)
	// // }
	// ps, err := repository.FindByAge(56)
	// // ps, err := repository.GetAll()
	// if err != nil {
	// 	log.Panic(err)
	// }
	// jsonData, err := json.Marshal(ps)
	// if err != nil {
	// 	log.Panic(err)
	// }
	// fmt.Println(string(jsonData))
	// os.Exit(0)

	exRepository, err := NewExcludeRepository("exclude_actors.txt")
	if err != nil {
		log.Panicf("failed to NewExcludeRepository: %v", err)
	}

	config := Config{Port: "8080"}

	service := NewActorService(config, repository, exRepository)

	server := NewServer(config, service)
	server.Run()
}

type Config struct {
	Port string
}

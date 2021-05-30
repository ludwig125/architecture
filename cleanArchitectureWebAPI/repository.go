package main

type ActorRepository interface {
	GetAll() ([]Actor, error)
	SearchByID(int) ([]Actor, error)
	SearchByName(string) ([]Actor, error)
	SearchByAge(int) ([]Actor, error)
	Update(Actor) error
	DeleteByID(int) error
}

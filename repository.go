package main

type ActorRepository interface {
	GetAll() ([]Actor, error)
	FindByID(int) ([]Actor, error)
	FindByName(string) ([]Actor, error)
	FindByAge(int) ([]Actor, error)
	Update(Actor) error
	DeleteByID(int) error
}

type ExcludeRepository interface {
	Excluded() ([]string, error)
}

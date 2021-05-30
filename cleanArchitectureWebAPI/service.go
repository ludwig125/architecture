package main

import (
	"fmt"
)

type ActorService interface {
	GetAll() ([]Actor, error)
	Search(RequestCond) ([]Actor, error)
	Update(Actor) error
	DeleteByID(int) error
}

type actorService struct {
	config     Config
	repository ActorRepository
}

// interfaceを実装しているか保証する
// See: http://golang.org/doc/faq#guarantee_satisfies_interface
var _ ActorService = (*actorService)(nil)

func NewActorService(config Config, repository ActorRepository) ActorService {
	return &actorService{config: config, repository: repository}
}

func (s *actorService) GetAll() ([]Actor, error) {
	return s.repository.GetAll()
}

func (s *actorService) Search(cond RequestCond) ([]Actor, error) {
	var as []Actor
	var err error

	switch {
	case cond.ID != 0:
		as, err = s.repository.SearchByID(cond.ID)
	case cond.Name != "":
		as, err = s.repository.SearchByName(cond.Name)
	case cond.Age != 0:
		as, err = s.repository.SearchByAge(cond.Age)
	default:
		return nil, fmt.Errorf("invalid condition: %#v", cond)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to search: %v", err)
	}
	if len(as) == 0 {
		return nil, fmt.Errorf("did not meet the conditions %#v", cond)
	}

	return as, nil
}

func (s *actorService) Update(a Actor) error {
	return s.repository.Update(a)
}

func (s *actorService) DeleteByID(id int) error {
	return s.repository.DeleteByID(id)
}

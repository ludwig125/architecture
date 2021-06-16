package main

import (
	"fmt"
	"log"
)

type ActorService interface {
	GetAll() ([]Actor, error)
	Find(Actor) ([]Actor, error)
	Update(Actor) error
	DeleteByID(int) error
	// Excluded() ([]string, error)
}

type actorService struct {
	config       Config
	repository   ActorRepository
	exRepository ExcludeRepository
}

// interfaceを実装しているか保証する
// See: http://golang.org/doc/faq#guarantee_satisfies_interface
var _ ActorService = (*actorService)(nil)

func NewActorService(config Config, repository ActorRepository, exRepository ExcludeRepository) ActorService {
	return &actorService{config: config, repository: repository, exRepository: exRepository}
}

func (s *actorService) GetAll() ([]Actor, error) {
	as, err := s.repository.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to GetAll: %v", err)
	}

	// 除外対象のものをのぞいて返す
	es, err := s.exRepository.Excluded()
	if err != nil {
		return nil, fmt.Errorf("failed to get Excluded: %v", err)
	}
	return fileterActors(as, es), nil
}

func (s *actorService) Find(cond Actor) ([]Actor, error) {
	as, err := s.findActors(cond)
	if err != nil {
		return nil, fmt.Errorf("failed to findActors: %v", err)
	}

	// 除外対象のものをのぞいて返す
	es, err := s.exRepository.Excluded()
	if err != nil {
		return nil, fmt.Errorf("failed to get Excluded: %v", err)
	}
	return fileterActors(as, es), nil
}

func (s *actorService) findActors(cond Actor) ([]Actor, error) {
	var as []Actor
	var err error

	switch {
	case cond.ID != 0:
		as, err = s.repository.FindByID(cond.ID)
	case cond.Name != "":
		as, err = s.repository.FindByName(cond.Name)
	case cond.Age != 0:
		as, err = s.repository.FindByAge(cond.Age)
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
	// すでにあるものについてはUpdateしない
	// updateとして指定されるものにはIDがないので、repositoryのINSERT OR REPLACE INTOでは制御できない
	as, err := s.repository.FindByName(a.Name)
	if err != nil {
		return fmt.Errorf("failed to FindByName: %v", err)
	}
	if len(as) > 0 {
		return fmt.Errorf("actor %s already exists", a.Name)
	}
	return s.repository.Update(a)
}

func (s *actorService) DeleteByID(id int) error {
	return s.repository.DeleteByID(id)
}

func fileterActors(as []Actor, es []string) []Actor {
	var after []Actor
	for _, a := range as {
		if sliceHas(es, a.Name) {
			log.Println("found excluded actor", a.Name)
			continue
		}
		after = append(after, a)
	}
	return after
}

func sliceHas(es []string, target string) bool {
	for _, e := range es {
		if e == target {
			return true
		}
	}
	return false
}

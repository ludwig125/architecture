package main

import (
	"errors"
	"fmt"
	"strconv"
)

// Actor is actor struct with json.
type Actor struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// RequestCond is request condition to search actor.
type RequestCond struct {
	ID   int
	Name string
	Age  int
}

func NewRequestCond(id, name, age string) (*RequestCond, error) {
	// 同時に複数パラメータを指定するとエラーにする
	multiErr := fmt.Errorf("multiple variables. id: %s, name: %s, age: %s, choose one variable", id, name, age)

	switch {
	case id != "":
		if name != "" || age != "" {
			return nil, multiErr
		}

		iID, err := strconv.Atoi(id)
		if err != nil {
			return nil, fmt.Errorf("failed to convert id '%s' to int", id)
		}
		return &RequestCond{
			ID:   iID,
			Name: "",
			Age:  0,
		}, nil
	case name != "":
		if id != "" || age != "" {
			return nil, multiErr
		}

		return &RequestCond{
			ID:   0,
			Name: name,
			Age:  0,
		}, nil
	case age != "":
		if id != "" || name != "" {
			return nil, multiErr
		}

		iAge, err := strconv.Atoi(age)
		if err != nil {
			return nil, fmt.Errorf("failed to convert age '%s' to int", age)
		}
		return &RequestCond{
			ID:   0,
			Name: "",
			Age:  iAge,
		}, nil
	default:
		return nil, errors.New("invalid")
	}
}

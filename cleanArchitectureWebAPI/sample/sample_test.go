package main

import (
	"fmt"
	"testing"
	"time"
)

func TestLogStreaming(t *testing.T) {
	for i := 0; i < 5; i++ {
		time.Sleep(300 * time.Millisecond)
		fmt.Println("fmt.Println:", i)
		t.Log("t.Log:", i)
	}
}

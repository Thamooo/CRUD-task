package main

import (
	"testing"
)

func Test_checkAge(t *testing.T) {
    result := checkAge("2000-01-01", 20, 90)
    if result != true {
       t.Errorf("Sum was incorrect, got: %d, want: %d.", result, true)
    }
}